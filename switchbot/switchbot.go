package switchbot

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

type SwitchBot struct {
	address bluetooth.Address
}

func ScanSwitchBot(ctx context.Context, address string) (*SwitchBot, error) {
	// Enable BLE interface.
	if err := adapter.Enable(); err != nil {
		return nil, fmt.Errorf("enable BLE stack: %w", err)
	}

	// Start scanning.
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	found := make(chan bluetooth.ScanResult, 1)
	errCh := make(chan error, 1)
	err := adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
		if strings.HasSuffix(strings.ToLower(strings.ReplaceAll(device.Address.String(), ":", "")), address) {
			if err := adapter.StopScan(); err != nil {
				errCh <- err
				return
			}
			found <- device
		}
	})
	if err != nil {
		return nil, fmt.Errorf("start scan: %w", err)
	}

	var result bluetooth.ScanResult
	select {
	case <-ctx.Done():
		if err := adapter.StopScan(); err != nil {
			return nil, fmt.Errorf("stop scan: %w", err)
		}
		return nil, fmt.Errorf("scan timeout")
	case err := <-errCh:
		return nil, fmt.Errorf("scan error: %w", err)
	case result = <-found:
		log.Println("found SwitchBot:", result.Address.String())
	}

	return &SwitchBot{address: result.Address}, nil
}

func (bot *SwitchBot) TurnOn(ctx context.Context) error {
	// 再接続のためのリトライロジック
	const maxRetries = 3
	var lastErr error

	for i := range maxRetries {
		log.Printf("Attempting to turn on SwitchBot (attempt %d/%d)", i+1, maxRetries)

		if err := bot.turnOnAttempt(ctx); err != nil {
			lastErr = err
			log.Printf("Turn on attempt %d failed: %v", i+1, err)

			// D-Busエラーの場合、アダプターを再初期化して再試行
			if strings.Contains(err.Error(), "org.freedesktop.DBus.Properties") ||
				strings.Contains(err.Error(), "doesn't exist") ||
				strings.Contains(err.Error(), "Method") && strings.Contains(err.Error(), "signature") {

				log.Printf("D-Bus error detected, reinitializing adapter (attempt %d/%d): %v", i+1, maxRetries, err)
				if reinitErr := adapter.Enable(); reinitErr != nil {
					log.Printf("Failed to reinitialize adapter: %v", reinitErr)
				} else {
					log.Printf("Successfully reinitialized adapter")
				}

				if i < maxRetries-1 {
					sleepDuration := time.Second * time.Duration(i+1) // 段階的バックオフ
					log.Printf("Waiting %v before retry", sleepDuration)
					time.Sleep(sleepDuration)
					continue
				}
			}

			// その他のエラーの場合は即座に失敗
			if i == 0 {
				return err
			}
		} else {
			log.Printf("Successfully turned on SwitchBot on attempt %d", i+1)
			return nil
		}
	}

	return fmt.Errorf("failed to turn on after %d retries, last error: %w", maxRetries, lastErr)
}

func (bot *SwitchBot) turnOnAttempt(ctx context.Context) error {
	device, err := adapter.Connect(bot.address, bluetooth.ConnectionParams{})
	if err != nil {
		return fmt.Errorf("connect to device: %w", err)
	}
	defer device.Disconnect()

	select {
	case <-ctx.Done():
		return fmt.Errorf("before discover services: %w", ctx.Err())
	default:
	}

	services, err := device.DiscoverServices([]bluetooth.UUID{mustParseUUID("cba20d00-224d-11e6-9fb8-0002a5d5c51b")})
	if err != nil {
		return fmt.Errorf("discover services: %w", err)
	}
	if len(services) == 0 {
		return fmt.Errorf("no bot service found")
	}

	botService := services[0]

	select {
	case <-ctx.Done():
		return fmt.Errorf("before discover characteristics: %w", ctx.Err())
	default:
	}

	chars, err := botService.DiscoverCharacteristics([]bluetooth.UUID{
		mustParseUUID("cba20003-224d-11e6-9fb8-0002a5d5c51b"),
		mustParseUUID("cba20002-224d-11e6-9fb8-0002a5d5c51b"),
	})
	if err != nil {
		return fmt.Errorf("discover characteristics: %w", err)
	}
	if len(chars) <= 1 {
		return fmt.Errorf("characteristic not found")
	}
	var notifyChar, writeChar bluetooth.DeviceCharacteristic
	if chars[0].UUID() == mustParseUUID("cba20003-224d-11e6-9fb8-0002a5d5c51b") {
		notifyChar = chars[0]
		writeChar = chars[1]
	} else {
		notifyChar = chars[1]
		writeChar = chars[0]
	}

	select {
	case <-ctx.Done():
		return fmt.Errorf("before enable notifications: %w", ctx.Err())
	default:
	}

	// Subscribe to notifications from the bot.
	if err := notifyChar.EnableNotifications(func(value []byte) {
		log.Printf("notification: %x\n", value)
	}); err != nil {
		return fmt.Errorf("enable notification: %w", err)
	}

	select {
	case <-ctx.Done():
		return fmt.Errorf("before write command: %w", ctx.Err())
	default:
	}

	// Write command to the characteristic to turn on the bot.
	cmd := []byte{0x57, 0x01, 0x00}
	_, err = writeChar.WriteWithoutResponse(cmd)
	if err != nil {
		return fmt.Errorf("write command: %w", err)
	}

	return nil
}

func mustParseUUID(s string) bluetooth.UUID {
	uuid, err := bluetooth.ParseUUID(s)
	if err != nil {
		panic("parse UUID: " + err.Error())
	}
	return uuid
}
