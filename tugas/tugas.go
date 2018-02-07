package tugas

import (
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/tokopedia/gosample/redis"
	db "github.com/tokopedia/gosample/utils/db"
)

type Order struct {
	OrderId       int           `json:"order_id"`
	InvoiceRefNum string        `json:"invoice_ref_num"`
	OrderStatus   int           `json:"order_status"`
	Shop          Shop          `json:"shop"`
	OrderDetail   []OrderDetail `json:"order_detail"`
}

type Shop struct {
	ShopId     int    `json:"shop_id"`
	ShopName   string `json:"shop_name"`
	ShopDomain string `json:"shop_domain"`
	ShopStatus uint8  `json:"shop_status"`
}

type Product struct {
	ProductId          int    `json:"product_id"`
	ProductName        string `json:"product_name"`
	ProductDescription string `json:"product_description"`
	ProductStatus      uint8  `json:"product_status"`
}

type OrderDetail struct {
	QuantityDelivered uint8   `json:"quantity_delivered"`
	Product           Product `json:"product"`
}

func HandleGetOrderAndOrderDetail(w http.ResponseWriter, r *http.Request) {

	sample_orders, err := getOrders()

	if err != nil {
		panic(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	var x []byte
	x, _ = json.Marshal(sample_orders)

	w.WriteHeader(http.StatusOK)
	w.Write(x)

	return
}

func HandleGetRedisExample(w http.ResponseWriter, r *http.Request) {
	data, _ := redis.Get("order_list:juna")

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(data))

	return
}

func getOrders() ([]Order, error) {
	rows, err := db.DB.DBTOrder.Query("SELECT order_id, invoice_ref_num, order_status, shop_id  FROM ws_order LIMIT 10")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var orders []Order
	for rows.Next() {
		var order_id int
		var invoice_ref_num string
		var order_status int
		var shop_id int
		if err := rows.Scan(&order_id, &invoice_ref_num, &order_status, &shop_id); err != nil {
			return orders, err
		}
		s, _ := getShop(shop_id)

		ods, _ := getOrderDetails(order_id)
		o := Order{
			OrderId:       order_id,
			InvoiceRefNum: invoice_ref_num,
			OrderStatus:   order_status,
			Shop:          s,
			OrderDetail:   ods,
		}

		orders = append(orders, o)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return orders, nil
}

func getOrderDetails(order_id int) ([]OrderDetail, error) {
	rows, err := db.DB.DBTOrder.Query("SELECT quantity_deliver, product_id  FROM ws_order_dtl WHERE order_id = $1", order_id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var result []OrderDetail

	for rows.Next() {
		var quantity_deliver uint8
		var product_id int

		if err := rows.Scan(&quantity_deliver, &product_id); err != nil {
			return result, err
		}

		temp_product, _ := getProduct(product_id)
		od := OrderDetail{
			Product:           temp_product,
			QuantityDelivered: quantity_deliver,
		}

		result = append(result, od)

	}
	return result, nil
}

func getProduct(product_id int) (Product, error) {
	rows, err := db.DB.DBTDev.Query("SELECT product_id, product_name, short_desc, status  FROM ws_product WHERE product_id = $1 LIMIT 1", product_id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var result Product

	for rows.Next() {
		var product_id int
		var product_name string
		var short_desc string
		var status uint8
		if err := rows.Scan(&product_id, &product_name, &short_desc, &status); err != nil {
			return result, err
		}
		result = Product{
			ProductId:          product_id,
			ProductName:        product_name,
			ProductDescription: short_desc,
			ProductStatus:      status,
		}

	}
	return result, nil
}

func getShop(shop_id int) (Shop, error) {
	rows, err := db.DB.DBTDev.Query(`SELECT shop_id, shop_name, domain, status FROM ws_shop WHERE shop_id = $1 LIMIT 1`, shop_id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var result Shop

	for rows.Next() {
		var shop_id int
		var shop_name string
		var domain string
		var status uint8
		if err := rows.Scan(&shop_id, &shop_name, &domain, &status); err != nil {
			return result, err
		}
		result = Shop{
			ShopId:     shop_id,
			ShopName:   shop_name,
			ShopDomain: domain,
			ShopStatus: status,
		}

	}
	return result, nil
}
