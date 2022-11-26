package main

//TODO: file containing IPs and mask
//		handle improper input
//		general clean up, DRY code
import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var REF = map[string][8]int{
	"1":        {1, 2, 3, 4, 5, 6, 7, 8},
	"2":        {9, 10, 11, 12, 13, 14, 15, 16},
	"3":        {17, 18, 19, 20, 21, 22, 23, 24},
	"4":        {25, 26, 27, 28, 29, 30},
	"magic":    {128, 64, 32, 16, 8, 4, 2, 1},
	"mask":     {128, 192, 224, 240, 248, 252, 254, 255},
	"networks": {2, 4, 8, 16, 32, 64, 128, 256},
}

// mag_num ==> hosts per subnet
var mag_num int

// Index ==> index for all look ups
var Index int

// global IP & subnet mask var
var IP string
var MASK string

func IndexOf(arr [8]int, item int) int {
	for i, v := range arr {
		if v == item {
			return i
		}
	}
	return -1
}
func getSubnet(ip, mask string) ([4]string, []string) {
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
			Index = IndexOf(mask_REF, mask_int)
			mag_num = REF["magic"][Index]
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
			//IPoctet number + magic num -1
			subnet_num, err := strconv.Atoi(subnet_arr[ix])
			handle_error(err)
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

func printIParrays(subnet_arr [4]string, broadcast_arr [4]string, host_1 [4]string, host_2 [4]string, net_num int) {
	subnet_str := strings.Join(subnet_arr[:], ".")
	broadcast_str := strings.Join(broadcast_arr[:], ".")
	host1_str := strings.Join(host_1[:], ".")
	host2_str := strings.Join(host_2[:], ".")

	fmt.Printf("\nIP Address: %v\nSubnet Mask: %v\n", IP, MASK)
	fmt.Printf("Subnet ID: %v\n", subnet_str)
	fmt.Printf("Broadcast Address: %v\n", broadcast_str)
	fmt.Printf("Host Range: %v --> %v\n", host1_str, host2_str)
	fmt.Printf("Number of Networks: %v\n", net_num)
	fmt.Printf("Number of Addresses/Subnet: %v\n", mag_num)
}

//func subMain() {
//	fmt.Println("Enter the IP & Subnet Mask")
//	var IPstring
//	var mask string
//	fmt.Scanln(&ip, &mask)
//	//fmt.Printf("IP: %-5v\nMask: %-5v", ip, mask)
//	subnet_arr, mask_arr := getSubnet(ip, mask)

//	broadcast_arr, subnet_arr := getBroadcast(mask_arr, subnet_arr)
//	//1st available host
//	host_1, host_2 := getHosts(subnet_arr, broadcast_arr)

//	// printIParrays(subnet_arr, broadcast_arr, host_1, host_2, net_num, mag_num)

//}

var HELP = "-h --> Help\n-c --> for slash notation AKA CIDR\n-m --> Subnet mask provided"

func arg_handler(args []string) []string {
	//func that takes all the args, parses them and spits out a
	//slice with 1st being whether or not its CIDR, or mask is given

	//get relevant
	var ip_mask_arr []string
	useful_args := args[1:]
	if useful_args[0] == "-c" {
		//get the ip/CIDR from useful_args
		str_res := useful_args[1]
		ip_cidr := strings.Split(str_res, "/")
		ip, cidr := ip_cidr[0], ip_cidr[1]
		ip_mask_arr = append(ip_mask_arr, "CIDR", ip, cidr)
	} else if useful_args[0] == "-m" {
		str_res := useful_args[1]
		ip_mask := strings.Split(str_res, "-")
		ip, mask := ip_mask[0], ip_mask[1]
		ip_mask_arr = append(ip_mask_arr, "mask", ip, mask)

	}
	return ip_mask_arr
}

func handle_error(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func toInt(str string) int {
	res, err := strconv.Atoi(str)
	handle_error(err)
	return res
}
func toStr(num int) string {
	res := strconv.Itoa(num)
	return res
}

var myErr error

func make_mask(ci string) (string, error) {
	//Func that masks a mask from CIDR number

	//convert to integer
	ci_int := toInt(ci)
	//loop over REF map
	//create mask by placing mask num at interesting octet
	errCheck := false
	var mask_str bytes.Buffer
	var myErr error
	for i, arr := range REF {
		arr_ind := IndexOf(arr, ci_int)
		if arr_ind != -1 {
			//go a match, now we get octet of interest and mask mask num
			octet_num := toInt(i)
			mask_num := REF["mask"][arr_ind]
			mask_oct := toStr(mask_num)
			// fmt.Println(octet_num, mask_num, mask_oct)

			flag := false // for if we hit the oct_num we now add 0
			for o := 1; o < 5; o++ {
				if o == octet_num {
					mask_str.WriteString(mask_oct)
					flag = true
				} else {
					if flag == true {
						mask_str.WriteString("0")
					} else {
						mask_str.WriteString("255")
					}
				}
				if o < 4 {
					mask_str.WriteString(".")
				}
			}
			break
		}
	}
	if errCheck != false {
		myErr = errors.New("Didnt find a number in REF. Should not happen!")
	}
	// fmt.Println(mask_str.String())
	return mask_str.String(), myErr
}

func process(ip, mask string) {
	// processes the IP and mask in string and encapsulates functions that
	// get the get the subnet ID, Broadcast address and Host range
	subnet_arr, mask_arr := getSubnet(ip, mask)

	broadcast_arr, subnet_arr := getBroadcast(mask_arr, subnet_arr)
	//1st available host
	host_1, host_2 := getHosts(subnet_arr, broadcast_arr)
	//get network number
	// log.Println(Index)
	network_num := REF["networks"][Index]
	// log.Println(IP, MASK)
	printIParrays(subnet_arr, broadcast_arr, host_1, host_2, network_num)
}

func cli() {
	args := os.Args
	//fmt.Println(len(args))
	len_args := len(args)
	switch {
	case len_args == 1:
		fmt.Printf("%v", HELP)
	default:
		// fmt.Println("Arguments given")
		// fmt.Printf("%v",args)
		user_data := arg_handler(args)
		// fmt.Println(user_data)
		mode := user_data[0]
		m_num := user_data[2]
		IP = user_data[1]
		if mode == "CIDR" {
			msk, err := make_mask(m_num)
			handle_error(err)
			MASK = msk
			log.Println(IP, MASK)
			process(IP, MASK)
		} else if mode == "mask" {
			MASK = user_data[2]
			process(IP, MASK)
		}
	}
}

func main() {

	//fmt.Println(REF)
	cli()
	// subMain()

}
