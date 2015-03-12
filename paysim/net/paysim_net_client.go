package net

import (
	"bytes"
	"encoding/hex"
	"errors"
	"github.com/rkbalgi/go/iso8583"
	p_nc "github.com/rkbalgi/go/net"
	"log"
)

var nc_map map[string]*p_nc.NetCatClient

func init() {
	nc_map = make(map[string]*p_nc.NetCatClient, 10)
}

//SendIsoMsg will send a msg over to the host
//receive and return the parsed response message
//TODO:: there is need to keep a map of the requests/responses
//and manage them via a key in case there are multiple messages
func SendIsoMsg(connection_str string,
	mli string,
	iso_msg *iso8583.Iso8583Message) (*iso8583.Iso8583Message, error) {

	var mli_type p_nc.MliType

	if mli == "2I" {
		mli_type = p_nc.MLI_2I
	} else {
		mli_type = p_nc.MLI_2E
	}

	nc, ok := nc_map[connection_str]
	//lets check if nc is still connected
	if ok && !nc.IsConnected(){
		ok=false;
	}
	
	if !ok {
		nc = p_nc.NewNetCatClient(connection_str, mli_type)
		err := nc.OpenConnection()
		if err != nil {
			return nil, err
		}
		log.Println("new tcp/ip connection opened to -", connection_str)
	}

	//we have a client  now
	req_msg_data := iso_msg.Bytes()
	log.Println("sending data \n", hex.Dump(req_msg_data), "\n")
	nc.Write(req_msg_data)
	resp_msg_data, err := nc.ReadNextPacket()
	if err != nil {
		return nil, err
	}
	log.Println("received data \n", hex.Dump(resp_msg_data), "\n")

	resp_iso_msg := iso8583.NewIso8583Message(iso_msg.SpecName())
	msg_buf := bytes.NewBuffer(resp_msg_data)
	err = resp_iso_msg.Parse(msg_buf)
	if err != nil {
		return nil, errors.New("error parsing response data")
	}

	return resp_iso_msg, nil

}
