package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"os"
)

var REF = map[string][8]int{
	"1": {1,2,3,4,5,6,7,8},
	"2": {9,10,11,12,13,14,15,16},
	"3": {17,18,19,20,21,22,23,24},
	"4": {25,26,27,28,29,30},
	"magic": {128, 64, 32, 16, 8, 4, 2, 1},
	"mask":  {128, 192, 224, 240, 248, 252, 254, 255},
}

func IndexOf(arr [8]int, item int) int {
	for i, v := range arr {
		if v == item {
			return i
		}
	}
	return -1
}
func getSubnet(ip string, mask string) ([4]string, []string) {
	ip_arr := strings.Split(ip, ".")
	mask_arr := strings.Split(mask, ".")

	//fmt.Printf("IP: %v\nMask: %v\n", ip_arr, mask_arr)

	//Subnet Id loop
	var subnet_arr [4]string
	for ix := 0; ix < 4; ix++ {
		switch {
		case mask_arr[ix] == "255":
			subnet_arr[ix] = ip_arr[ix]
		case mask_arr[ix] == "0":
			subnet_arr[ix] = "0"
		default:
			//get the index of the mask number
			mask_int, err := strconv.Atoi(mask_arr[ix])
			if err != nil {
				log.Fatal(err)
			}
			mask_REF := REF["mask"]
			index := IndexOf(mask_REF, mask_int)
			mag_num := REF["magic"][index]
			collect := 0
			ip_int, err := strconv.Atoi(ip_arr[ix])
			if err != nil {
				log.Fatal(err)
			}
			for m := 1; m < 255; m++ {
				result := m * mag_num
				if result <= ip_int {
					collect = result
				} else {
					break
				}
			}
			subnet_arr[ix] = strconv.Itoa(collect)
		}
	}
	return subnet_arr, mask_arr
}

func getBroadcast(mask_arr []string, subnet_arr [4]string) ([4]string, [4]string) {
	//Broadcast address
	var broadcast_arr [4]string
	for ix := 0; ix < 4; ix++ {
		switch {
		case mask_arr[ix] == "255":
			broadcast_arr[ix] = subnet_arr[ix]
		case mask_arr[ix] == "0":
			broadcast_arr[ix] = "255"
		default:
			//boilerplate
			mask_int, err := strconv.Atoi(mask_arr[ix])
			if err != nil {
				log.Fatal(err)
			}
			mask_REF := REF["mask"]
			index := IndexOf(mask_REF, mask_int)
			mag_num := REF["magic"][index]
			//ip octet number + magic num -1
			subnet_num, err := strconv.Atoi(subnet_arr[ix])
			if err != nil {
				log.Fatal(err)
			}
			result := subnet_num + mag_num - 1
			str_res := strconv.Itoa(result)
			broadcast_arr[ix] = str_res
		}
	}
	return broadcast_arr, subnet_arr
}

func getHosts(subnet_arr [4]string, broadcast_arr [4]string) ([4]string, [4]string) {
	host_1 := subnet_arr
	host_1[3] = "1"

	// last Host
	host_2 := broadcast_arr
	host_2[3] = "254"

	return host_1, host_2
}

func printIParrays(subnet_arr [4]string, broadcast_arr [4]string, host_1 [4]string, host_2 [4]string) {
	subnet_str := strings.Join(subnet_arr[:], ".")
	broadcast_str := strings.Join(broadcast_arr[:], ".")
	host1_str := strings.Join(host_1[:], ".")
	host2_str := strings.Join(host_2[:], ".")

	fmt.Printf("\nSubnet ID: %v\n", subnet_str)
	fmt.Printf("Broadcast Address: %v\n", broadcast_str)
	fmt.Printf("Host Range: %v --> %v\n", host1_str, host2_str)
}

func subMain(){
	fmt.Println("Enter the IP & Subnet Mask")
	var ip string
	var mask string
	fmt.Scanln(&ip, &mask)
	//fmt.Printf("IP: %-5v\nMask: %-5v", ip, mask)
	subnet_arr, mask_arr := getSubnet(ip, mask)

	broadcast_arr, subnet_arr := getBroadcast(mask_arr, subnet_arr)
	//1st available host
	host_1, host_2 := getHosts(subnet_arr, broadcast_arr)

	printIParrays(subnet_arr, broadcast_arr, host_1, host_2)

}

var HELP = "-h --> Help\n-c --> for slash notation AKA CIDR\n-m --> Subnet mask provided"
func cli(){
	args := os.Args
	//fmt.Println(len(args))
	len_args := len(args)
	switch{
	case len_args == 1:
		fmt.Printf("%v",HELP)
	default:
		fmt.Println("Arguments given")
		fmt.Printf("%v",args)
	}
}

func main() {

	//fmt.Println(REF)
	cli()
	subMain()

}
