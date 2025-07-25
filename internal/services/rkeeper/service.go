package rkeeper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/AlexChe360/procash/internal/config"
)

// Базовый POST-запрос
func post(cfg config.Config, taskType string, params map[string]any) (map[string]any, error) {
	payload := map[string]any{
		"taskType": taskType,
		"params":   params,
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", cfg.RKeeperBaseURL, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("AggregatorAuthentication", "Token "+cfg.RKeeperToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	var result map[string]any
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Получение tableCode по externalNumber
func GetTableCode(cfg config.Config, restaurantID int, tableNumber string) (int, error) {
	params := map[string]any{
		"sync": map[string]any{
			"objectId": restaurantID,
			"timeout":  120,
		},
	}

	resp, err := post(cfg, "GetTableList", params)
	if err != nil {
		return 0, err
	}

	tablesRaw, ok := resp["tables"]
	if !ok {
		log.Println("❌ Поле 'tables' отсутствует в ответе от rKeeper")
		return 0, fmt.Errorf("missing 'tables' field")
	}

	tables, ok := tablesRaw.([]any)
	if !ok {
		log.Println("❌ Поле 'tables' не массив")
		return 0, fmt.Errorf("'tables' is not a slice")
	}

	for _, t := range tables {
		table, ok := t.(map[string]any)
		if !ok {
			continue // или log.Println("⚠️ Невалидная структура table")
		}

		if table["externalNumber"] == tableNumber {
			codeFloat, ok := table["code"].(float64)
			if !ok {
				log.Println("⚠️ Невалидный тип 'code'")
				continue
			}
			return int(codeFloat), nil
		}
	}

	return 0, fmt.Errorf("table %s not found", tableNumber)
}

// Получение GUID заказа и ID официанта
func GetOrderInfo(cfg config.Config, restaurantID int, tableCode int) (string, string, error) {
	params := map[string]any{
		"sync": map[string]any{
			"objectId": restaurantID,
			"timeout":  120,
		},
		"tableCode":  tableCode,
		"withClosed": false,
	}

	resp, err := post(cfg, "GetOrderList", params)
	if err != nil {
		return "", "", err
	}

	ordersRaw, ok := resp["orders"]
	if !ok {
		return "", "", fmt.Errorf("❌ поле 'orders' отсутствует в ответе от rKeeper")
	}

	orders, ok := ordersRaw.([]any)
	if !ok {
		return "", "", fmt.Errorf("❌ неверный формат 'orders' в ответе от rKeeper")
	}

	if len(orders) == 0 {
		return "", "", fmt.Errorf("❌ нет открытых заказов")
	}

	orderMap, ok := orders[0].(map[string]any)
	if !ok {
		return "", "", fmt.Errorf("❌ неверный формат заказа")
	}

	guid, ok1 := orderMap["guid"].(string)
	waiterId, ok2 := orderMap["waiterId"].(string)
	if !ok1 || !ok2 {
		return "", "", fmt.Errorf("❌ отсутствуют поля 'guid' или 'waiterId'")
	}

	return guid, waiterId, nil
}

// Детали заказа (товары и сумма)
func GetOrderDetails(cfg config.Config, restaurantID int, orderGUID string) (items []map[string]any, totalSum int, err error) {
	params := map[string]any{
		"sync": map[string]any{
			"objectId": restaurantID,
			"timeout":  120,
		},
		"orderGuid": orderGUID,
	}

	resp, err := post(cfg, "GetOrder", params)
	if err != nil {
		return nil, 0, err
	}

	rawItems, ok := resp["items"]
	if !ok {
		log.Println("❌ 'items' отсутствует в ответе от rKeeper")
		return nil, 0, fmt.Errorf("missing 'items' field")
	}

	itemsSlice, ok := rawItems.([]any)
	if !ok {
		log.Println("❌ 'items' не является массивом")
		return nil, 0, fmt.Errorf("'items' is not a slice")
	}

	for _, i := range itemsSlice {
		item, ok := i.(map[string]any)
		if !ok {
			log.Println("⚠️ Пропущен элемент: некорректная структура item")
			continue
		}
		items = append(items, map[string]any{
			"name":     item["name"],
			"quantity": item["quantity"],
			"amount":   item["amount"],
		})
	}

	rawSum, ok := resp["totalSum"]
	if !ok {
		log.Println("❌ 'totalSum' отсутствует в ответе")
		return items, 0, fmt.Errorf("missing 'totalSum' field")
	}

	sumFloat, ok := rawSum.(float64)
	if !ok {
		fmt.Println("❌ 'totalSum' не является числом")
		return items, 0, fmt.Errorf("invalid 'totalSum' type")
	}

	return items, int(sumFloat), nil
}

// Получение имени официанта
func GetWaiterName(cfg config.Config, restaurantID int, waiterID string) (string, error) {
	params := map[string]any{
		"sync": map[string]any{
			"objectId": restaurantID,
			"timeout":  120,
		},
	}

	resp, err := post(cfg, "GetEmployees", params)
	if err != nil {
		return "", err
	}

	employeesRaw, ok := resp["employees"]
	if !ok {
		fmt.Println("❌ 'employees' отсутствует в ответе от rKeeper")
		return "Unknown", fmt.Errorf("missing 'employees' field")
	}

	employees, ok := employeesRaw.([]any)
	if !ok {
		fmt.Println("❌ 'employees' не массив")
		return "Unknown", fmt.Errorf("employees' is not a slice")
	}

	for _, e := range employees {
		emp, ok := e.(map[string]any)
		if !ok {
			log.Println("⚠️ Пропущен элемент: некорректная структура employee")
			continue
		}
		id, ok := emp["id"].(string)
		if !ok {
			log.Println("⚠️ Пропущен элемент: 'id' не строка")
			continue
		}
		if id == waiterID {
			name, ok := emp["name"].(string)
			if !ok {
				log.Println("⚠️ 'name' не строка")
				return "Unknown", nil
			}
			return name, nil
		}
	}

	return "Unknown", nil
}
