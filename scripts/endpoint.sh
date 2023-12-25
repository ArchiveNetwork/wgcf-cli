#!/bin/bash

get_ips() {
    cleanup() {
        rm $ip6s_temp $ips_temp -f 
        exit 1
    }
    trap cleanup INT
    ##
    ips() {
        curl -sL https://gitlab.com/fscarmen/warp/-/raw/main/endpoint/ipv4 | shuf -n $IP_NUMBER > $ips_temp
        ipv4_array=()
        while IFS= read -r line; do
            if [[ $line =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
                ipv4_array+=("$line")
            fi
        done < $ips_temp
    }
    ip6s() {
        curl -sL https://gitlab.com/fscarmen/warp/-/raw/main/endpoint/ipv6 | shuf -n $IP_NUMBER > $ip6s_temp
        ipv6_array=()
        while IFS= read -r line; do
            if [[ $line =~ ^\[?([0-9a-fA-F:]+)\]?$ ]]; then
                ipv6_array+=("${BASH_REMATCH[1]}") 
            fi
        done < $ip6s_temp
    }
    ips_temp=$(mktemp)
    ip6s_temp=$(mktemp)

    if [[ -z $IPv ]];then
        IP_NUMBER=$((IP_NUMBER / 2))
        ips;ip6s
    elif [[ $IPv == 6 ]];then
        ip6s
        v6=6
    elif [[ $IPv == 4 ]];then
        ips
        v6=
    fi

    rm $ip6s_temp $ips_temp -f 
}

ping_ips() {
    cleanup() {
        rm $result_temp -f
        exit 1
    }
    trap cleanup INT
    ##
    echo "Testing. Don't be too anxious"

    for ip in "${ipv6_array[@]}" "${ipv4_array[@]}"; do
        result=$(ping -c 1 -W 0.1 "$ip")
        response_time=$(echo "$result" | grep -oP "time=\K[0-9.]+")
        if [[ ! -z "$response_time" ]]; then
            ping_results["$ip"]=$response_time
            echo "$ip ${ping_results["$ip"]} ms"
        fi
    done

    result_temp=$(mktemp)
    for ip in "${!ping_results[@]}"; do
        response_time="${ping_results[$ip]}"
        echo "$ip $response_time ms" >> $result_temp
    done

    if [[ $STORE == true ]];then 
        result=./result.txt
    else
        result=$(mktemp)
    fi

    sort -n -k 2 $result_temp -o $result

    rm $result_temp -f
}

replace(){
    cleanup() {
        wg-quick down wgcf
        exit 1
    }
    trap cleanup INT
    ##
    TEST=1
        TEST() {

            ports=2408 #(854 859 864 878 880 890 891 894 903 908 928 934 939 942 943 945 946 955 968 987 988 1002 1010 1014 1018 1070 1074 1180 1387 1843 2371 2408 2506 3138 3476 3581 3854 4177 4198 4233 5279 5956 7103 7152 7156 7281 7559 8319 8742 8854 8886)

            if cat $result | head -n 1| grep -oE '[0-9a-fA-F:]+:[0-9a-fA-F]+';then 
                endpoint=$(cat $result | head -n 1| grep -oE '[0-9a-fA-F:]+:[0-9a-fA-F]+')
                v6=6
            elif cat $result | head -n 1 | grep -Eo '\b([0-9]{1,3}\.){3}[0-9]{1,3}\b';then
                endpoint=$(cat $result | head -n 1 | grep -Eo '\b([0-9]{1,3}\.){3}[0-9]{1,3}\b')
                v6=
            else
                endpoint=
            fi
            if [[ -z $endpoint ]];then
                echo "Finished, no avilable endpoint found"
                rm $result
                wg-quick down wgcf >/dev/null
                return 1
            fi
            for port in ${ports[@]};do
                if [[ $v6 == 6 ]];then 
                    sed -i "s/Endpoint = .*:.*/Endpoint = [$endpoint]:$port/" /etc/wireguard/wgcf.conf
                else
                    sed -i "s/Endpoint = .*:.*/Endpoint = $endpoint:$port/" /etc/wireguard/wgcf.conf
                fi
                sed -i "s/ip$v6 saddr .* udp/ip$v6 saddr $endpoint udp/" /etc/wireguard/wgcf.nft.conf
                sed -i "s/ip$v6 daddr .* udp/ip$v6 daddr $endpoint udp/" /etc/wireguard/wgcf.nft.conf
                sed -i "s/sport \([0-9]\+\)/sport $port/g" /etc/wireguard/wgcf.nft.conf
                sed -i "s/dport \([0-9]\+\)/dport $port/g" /etc/wireguard/wgcf.nft.conf
                wg-quick down wgcf ;wg-quick up wgcf;ping -I wgcf -W 0.5 1.1.1.1 -c 1 > /dev/null
                sleep 0.5 &
                wait $! ; echo
                if [[ $v6 == 6 ]];then 
		            echo \[$endpoint\]:$port
                else
                    echo $endpoint:$port
                fi
                if (response=$(curl --Interface wgcf --max-time 0.1 -w '%{time_total}' -o /dev/null -s 1.1.1.1 2>&1); [ $? -ne 28 ] && (echo "$response"| awk '{printf "%.0f ms\n", $1*1000}')) || ping 1.1.1.1 -c 1 -I wgcf -W 0.1;then
                    FINDED=true
                    echo "Tested times:$TEST"
                else
                    FINDED=false
                    echo "Tested times:$TEST"
                fi
                [[ $FINDED == true ]] && break
                ((TEST++))
            done
            [[ $FINDED == false ]] && sed -ie "/$endpoint/d" $result
            [[ $FINDED == false ]] && TEST
        }
    TEST || main "$args"
}

menu() {
    ask_ip(){
        if [[ -z $IPv ]];then
            echo -en "\
Chose your IP version
1. IPv4
2. IPv6
3. Both  (Deafult)
Q. Exit
"
            read -p "" choice
            if [ -z "$choice" ]; then
            IPv=
            else
                case $choice in
                    1)
                        IPv=4
                        ;;
                    2)
                        IPv=6
                        ;;
                    3)
                       IPv=
                        ;;
                    q|Q)
                       exit 0
                        ;;
                    *)
                        echo "Invailed choice."
                        ask_ip
                    ;;
                esac
            fi
        fi
    }
    ask_store(){
        if [[ -z $STORE ]];then
            if [[ -f ./result.txt ]];then 
                read -p "\
\"./result.txt\" Found, do you wan't to use it ? (Y|n)\
" choice
                if [ -z "$choice" ]; then
                    STORE=true
                    result=./result.txt
                fi
            else 
                read -p "\
Do you wan't to store result to \"./result.txt\" ? (y|N)\
" choice
                if [ -z "$choice" ]; then
                    STORE=false
                fi
            fi  

            case $choice in
                y|Y)
                    STORE=true
                    result=./result.txt
                    ;;
                n|N)
                    STORE=false
                    ;;
                q|Q)
                   exit 0
                    ;;
                *)
                    if ! [ -z "$choice" ]; then
                        echo "Invailed choice."
                        ask_store
                    fi
                    ;;
            esac
        fi
    }
    ask_ip
    ask_store
}

judgement() {
    while [[ "$#" -gt '0' ]]; do
        case $1 in
            --ipv6|--v6|-6)
                IPv=6
                ;;
            --ipv4|--v4|-4)
                IPv=4
                ;;
            -s|--store)
                STORE=true
                result=./result.txt
                ;;
            -t|--time)
                [[ -z $2 ]] && echo "Please provide a valid number"&& exit 1
                IP_NUMBER="$2"
                shift
                ;;
            -y|--yes)
                ByPass=true
                ;;
            *)
                echo "Invalid argument: $arg"
                exit 1
                ;;
        esac
        shift
    done
}

main() {
    declare -A ping_results
    args="$@"
    judgement "$@"
    [[ -z $IP_NUMBER ]] && IP_NUMBER=50
    if ! [[ $ByPass == true ]];then
        menu
    fi
    if ! [[ -f ./result.txt ]] && [[ $STORE == true ]];then
        get_ips
        ping_ips  
    elif [[ $STORE == false ]];then
        get_ips
        ping_ips
    fi
    replace 
}

main