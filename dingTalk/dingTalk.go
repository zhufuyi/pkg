package dingTalk

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/CatchZeng/dingtalk"
)

const (
	limitTimeRange = 60 // 统计时间范围，单位为秒
	limitFrequency = 19 // 在统计时间范围内最多19次限制
)

//TokenSecret 钉钉机器人的信息
type TokenSecret struct {
	Name   string // 机器人名称
	Token  string // token
	Secret string // 密钥
}

// Client 钉钉客户端
type Client struct {
	currentIndex int                             // 当前使用的索引
	names        []string                        // 机器人名称
	clients      map[string]*dingtalk.Client     // 机器人名称对应的客户端
	timestamps   map[string]*[limitFrequency]int // 机器人对应的发送时刻表，用来判断20次/分钟限制
	mutex        *sync.Mutex                     // 互斥锁
}

var client *Client

// Init 初始化钉钉客户端
func Init(tss []TokenSecret) int {
	var (
		robotNames   []string
		robotClients = make(map[string]*dingtalk.Client)
		tts          = make(map[string]*[limitFrequency]int)
		count        = 0

		// 去重map
		l                   = len(tss)
		checkDuplicateName  = make(map[string]struct{}, l)
		checkDuplicateToken = make(map[string]struct{}, l)
	)

	for _, ts := range tss {
		if _, ok := checkDuplicateName[ts.Name]; ok {
			continue
		} else {
			checkDuplicateName[ts.Name] = struct{}{}
		}
		if _, ok := checkDuplicateToken[ts.Token]; ok {
			continue
		} else {
			checkDuplicateToken[ts.Token] = struct{}{}
		}

		robotNames = append(robotNames, ts.Name)
		robotClients[ts.Name] = dingtalk.NewClient(ts.Token, ts.Secret)
		tts[ts.Name] = new([limitFrequency]int)
		count++
	}

	client = &Client{
		currentIndex: -1, // 初始值
		names:        robotNames,
		clients:      robotClients,
		timestamps:   tts,
		mutex:        &sync.Mutex{},
	}

	return count
}

// Get 顺序获取钉钉机器人客户端client，并发安全
func Get() (*dingtalk.Client, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}
	client.mutex.Lock()
	defer client.mutex.Unlock()

	maxIndex := len(client.names) - 1
	if maxIndex < 0 {
		return nil, errors.New("dingding robot is empty")
	}

	// 索引循环
	if client.currentIndex < maxIndex {
		client.currentIndex++
	} else {
		client.currentIndex = 0
	}

	// 获取机器人名称
	name := client.names[client.currentIndex]

	// 判断是否达到发送频率限制(例如20次/分钟)
	timestamps := client.timestamps[name]
	wait, isTooFast := checkFrequency(timestamps, int(time.Now().Unix()))
	if isTooFast {
		return nil, fmt.Errorf("%s sending too fast, need wait %ds", name, wait)
	}

	// 获取机器人对象
	dtClient := client.clients[name]

	return dtClient, nil
}

// 判断是否超过频率限制，具体思路：在固定长度为n的数组中存放有序的时间戳环，随着时间推移，如果时间戳数量超过了
// 数组长度，就循环回来写入数组第一元素，永远只存放最新的n个时间戳值，而且是有序的，从时间戳环中判断最大最小时间戳值之差是否小于规定限制时间，
// 就可以判断发送速度超过频率是否超过限制频率，返回结果，如果为true，说明超过限制速率，并且返回具体需要等待时间(单位秒)
func checkFrequency(timestamps *[limitFrequency]int, nowSecond int) (int, bool) {
	minIndex := 0                // 最小值对应的索引
	lastTimeVal := timestamps[0] // 初始值
	waitTime := 0                // 下一次发送需要等待时间
	l := len(timestamps)

	for i, currentTimestamp := range timestamps {
		// 如果还有值为0情况，表示还没有达到速率限制
		if currentTimestamp == 0 {
			timestamps[i] = nowSecond
			return 0, false
		} else { // 所有值都以填充情况，获取最小值对应的index
			minIndex = i
			if lastTimeVal > currentTimestamp {
				break
			} else {
				if i == l-1 {
					minIndex = 0
				}
				lastTimeVal = currentTimestamp
			}
		}
	}

	flag := false
	val := nowSecond - timestamps[minIndex]
	if val < limitTimeRange { // 判断是否达到发送速率限制
		waitTime = limitTimeRange - val
		flag = true
	}
	timestamps[minIndex] = nowSecond

	return waitTime, flag
}

// ---------------------------------------------------------------------------------------

// NewClient 初始化client
func NewClient(token string, secret string) *dingtalk.Client {
	return dingtalk.NewClient(token, secret)
}

// NewTextMessage 实例化text消息类型对象
func NewTextMessage() *dingtalk.TextMessage {
	return &dingtalk.TextMessage{}
}

// NewMarkdownMessage 实例化markdown消息类型对象
func NewMarkdownMessage() *dingtalk.MarkdownMessage {
	return &dingtalk.MarkdownMessage{}
}

// NewLinkMessage 实例化link消息类型对象
func NewLinkMessage() *dingtalk.LinkMessage {
	return &dingtalk.LinkMessage{}
}

// NewFeedCardMessage 实例化feedCard消息类型对象
func NewFeedCardMessage() *dingtalk.FeedCardMessage {
	return &dingtalk.FeedCardMessage{}
}

// NewActionCardMessage 实例化actionCard消息类型对象
func NewActionCardMessage() *dingtalk.ActionCardMessage {
	return &dingtalk.ActionCardMessage{}
}
