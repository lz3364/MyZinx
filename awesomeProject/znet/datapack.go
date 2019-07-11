package znet

import (
	"awesomeProject/utils"
	"awesomeProject/ziface"
	"bytes"
	"encoding/binary"
	"errors"
)

type DataPack struct {
}

func (dp *DataPack) GetHeadLen() uint32 {
	return 8
}

func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {

	//创建一个存放bytes字节的缓冲
	// bytes类似于io流,实现了io操作的接口
	databuff := bytes.NewBuffer([]byte{})

	//写数据head部分的消息长度 DataLen
	if err := binary.Write(databuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	//写数据head部分的消息ID MsgId
	if err := binary.Write(databuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	//写数据body部分的消息 Data
	if err := binary.Write(databuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return databuff.Bytes(), nil
}

func (dp *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {

	dataBuff := bytes.NewBuffer(binaryData)
	msg := &Message{}

	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	//fmt.Println("[DEBUG] msg.DataLen:", msg.GetDataLen())

	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}
	//fmt.Println("[DEBUG] msg.GetMsgId:", msg.GetMsgId())

	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Data); err != nil {
		return nil, err
	}
	//fmt.Println("[DEBUG] msg.GetData:", string(msg.GetData()))

	if utils.GlobalObject.MaxPacketSize > 0 && msg.DataLen > utils.GlobalObject.MaxPacketSize {
		return nil, errors.New("Too large msg data recieved")
	}

	return msg, nil
}

func NewDataPack() *DataPack {
	return &DataPack{}
}
