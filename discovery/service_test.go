/**
 * @Author: lenovo
 * @Description:
 * @File:  service_test
 * @Version: 1.0.0
 * @Date: 2023/06/12 19:59
 */

package discovery

import "testing"

func TestServiceOrder1(t *testing.T) {
	ServiceOrder("localhost:8080")
}
func TestServiceOrder2(t *testing.T) {
	ServiceOrder("localhost:8081")
}
func TestServiceOrder3(t *testing.T) {
	ServiceOrder("localhost:8082")
}
