package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

func IndexOf(arr [8]int, item int) int {
	for i, v := range arr {
		if v == item {
			return i
		}
	}
	return -1
}

func main() {
	REF := map[string][8]int{
		"magic": {128, 64, 32, 16, 8, 4, 2, 1},
		"mask":  {128, 192, 224, 240, 248, 252, 254, 255},
	}
	//fmt.Println(REF)

	fmt.Println("Enter the IP & Subnet Mask")
	var ip string
	var mask string
	fmt.Scanln(&ip, &mask)
	//fmt.Printf("IP: %-5v\nMask: %-5v", ip, mask)

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

	//1st available host
	host_1 := subnet_arr
	host_1[3] = "1"

	// last Host
	host_2 := broadcast_arr
	host_2[3] = "254"

	subnet_str := strings.Join(subnet_arr[:], ".")
	broadcast_str := strings.Join(broadcast_arr[:], ".")
	host1_str := strings.Join(host_1[:], ".")
	host2_str := strings.Join(host_2[:], ".")

	fmt.Printf("\nSubnet ID: %v\n", subnet_str)
	fmt.Printf("Broadcast Address: %v\n", broadcast_str)
	fmt.Printf("Host Range: %v --> %v\n", host1_str, host2_str)

}