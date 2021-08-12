package ems

/*
#cgo CFLAGS: -I /opt/tibco/ems/8.4/include/tibems
#cgo LDFLAGS: -L/opt/tibco/ems/8.4/lib -L/opt/tibco/ems/8.4/lib/64 -ltibems64
#include <tibems.h>
#include <confact.h>
#include <emsadmin.h>
#include <dest.h>

*/
import "C"

type JMS struct {
	Status string
}
type jmsConnection struct {
	Status       string
	C_connection C.tibemsConnection
}
type ConnectionFactory struct {
}

type jmsSession struct {
	C_session    C.tibemsSession
	ErrorContext *C.tibemsErrorContext
}
type jmsMessageProducer struct {
	C_msgProducer C.tibemsMsgProducer
	C_destination C.tibemsDestination
	ErrorContext  *C.tibemsErrorContext
}
type jmsMessageConsumer struct {
	C_msgConsumer C.tibemsMsgConsumer
	Listener      MessageListener
	Callback      C.tibemsMsgCallback
	ErrorContext  *C.tibemsErrorContext
}
type jmsTextMessage struct {
	C_txtMsg     C.tibemsTextMsg
	ErrorContext *C.tibemsErrorContext
} //Connection   ss

type jmsMessage struct {
	C_msg        C.tibemsMsg
	text         string
	headers      map[string]string
	ErrorContext *C.tibemsErrorContext
}

type Connection interface {
	CreateSession(transacted bool, ack int) Session
	Start()
	Stop()
	Close()
}

//ConnectionFactory   dw
type ConnectionFactoryInterface interface {
	CreateConnection(emsServerURL string, username string, password string) Connection
}
type Session interface {
	CreateQueueProducer(queuename string) MessageProducer
	CreateTopicProducer(topicname string) MessageProducer
	CreateQueueConsumer(queuename string) MessageConsumer
	CreateTopicConsumer(queuename string) MessageConsumer
	Close()
}
type MessageProducer interface {
	Send(txtMsg TextMessage)
	Close()
}
type MessageConsumer interface {
	Receive() string
	SetMsgListener(MessageListener)
}
type TextMessage interface {
	Destroy()
}
type Message interface {
	GetText() string
	// SetText(text string)
	// GetHeaders() map[string]string
	// SetHeaders(headers map[string]string)
}
type MessageListener func(message Message)
