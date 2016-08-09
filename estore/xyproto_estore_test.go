package estore

import (
	"github.com/stretchr/testify/assert"
	es "github.com/sunrongya/eventsourcing"
	"github.com/xyproto/simplebolt"
	"os"
	"path"
	"strconv"
	"testing"
)

type mockEncoder struct {
	_event   []es.Event
	_version []int
	_out     string
	_err     error
}

func (this *mockEncoder) Encode(event es.Event, version int) (string, error) {
	this._event = append(this._event, event)
	this._version = append(this._version, version)
	return this._out + strconv.Itoa(version), this._err
}

type mockDecoder struct {
	_src        []string
	_exists     map[string]bool
	_outEvent   es.Event
	_outVersion int
	_errOut     error
}

func (this *mockDecoder) Decode(src string) (event es.Event, version int, err error) {
	if this._exists == nil {
		this._exists = make(map[string]bool)
	}
	this._src = append(this._src, src)
	if _, ok := this._exists[src]; !ok {
		this._exists[src] = true
		this._outVersion += 1
	}
	return this._outEvent, this._outVersion, this._errOut
}

func TestNewXyprotoEStore(t *testing.T) {
	db, err := simplebolt.New(path.Join(os.TempDir(), "bolt.db"))
	assert.NoError(t, err, "创建boltdb错误")
	defer db.Close()

	creator, factory := simplebolt.NewCreator(db), es.NewEventFactory()
	encoder, decoder := NewEncoder(factory), NewDecoder(factory)
	eventStore := NewXyprotoEStore(creator, encoder, decoder)

	assert.NotNil(t, eventStore, "eventStore应该不为nil")
	assert.Equal(t, creator, eventStore._creator, "creator应该相等")
	assert.Equal(t, encoder, eventStore._encoder, "encoder应该相等")
	assert.Equal(t, decoder, eventStore._decoder, "decoder应该相等")
	assert.NotNil(t, eventStore.EventChannel, "EventChannel应该不为nil")
}

func TestStoreEventForMock(t *testing.T) {
	db, _ := simplebolt.New(path.Join(os.TempDir(), "bolt.db"))
	defer db.Close()
	guid := es.NewGuid()
	creator := simplebolt.NewCreator(db)
	defer func() {
		list, _ := creator.NewList("Event" + string(guid))
		list.Clear()
		list.Remove()
	}()
	encoder := &mockEncoder{_out: "Event"}
	decoder := &mockDecoder{}
	events := []es.Event{
		&mockEvent1{es.WithGuid{guid}, "event1"},
		&mockEvent2{es.WithGuid{guid}, "event2"},
	}
	eventStore := NewXyprotoEStore(creator, encoder, decoder)
	err := eventStore.Update(guid, 0, events)

	assert.NoError(t, err, "添加正确事件不应该返回错误")
	assert.Equal(t, events, encoder._event, "事件解析序列有误")
	assert.Equal(t, []int{1, 2}, encoder._version, "事件解析版本序列有误")

	out, version := eventStore.Find(guid)
	assert.Equal(t, 2, version, "添加两条记录应该返回2")
	assert.Equal(t, []es.Event{decoder._outEvent, decoder._outEvent}, out, "返回的事件序列有误")
	assert.Equal(t, []string{"Event1", "Event2"}, decoder._src, "读出的数据有误")

	lastEvent, lastVersion, err := eventStore.getLastEvent(guid)
	assert.NoError(t, err, "读取最后的事件应该是没错误")
	assert.Equal(t, 2, lastVersion, "最新的version应该为2")
	assert.Equal(t, lastEvent, decoder._outEvent, "返回的event有误")
	assert.Equal(t, "Event2", decoder._src[2], "从db中返回的数据应为Event2")

	allEvents := eventStore.GetEvents(guid, 0, 8)
	assert.Equal(t, 2, len(allEvents), "返回的事件数应该为2")
	assert.Equal(t, []string{"Event1", "Event2", "Event2", "Event2", "Event1", "Event2"}, decoder._src)

	allEvents = eventStore.GetEvents(guid, 0, 1)
	assert.Equal(t, 1, len(allEvents), "返回的事件数应该为1")
	assert.Equal(t, []string{"Event1", "Event2", "Event2", "Event2", "Event1", "Event2", "Event2", "Event2"}, decoder._src)
}

func TestStoreEvents(t *testing.T) {
	db, _ := simplebolt.New(path.Join(os.TempDir(), "bolt.db"))
	defer db.Close()
	guid1, guid2 := es.NewGuid(), es.NewGuid()
	creator := simplebolt.NewCreator(db)
	defer func() {
		for _, guid := range []es.Guid{guid1, guid2} {
			list, _ := creator.NewList("Event" + string(guid))
			list.Clear()
			list.Remove()
		}
	}()
	eventFactory := es.NewEventFactory()
	eventFactory.RegisterAggregate(new(mockAggregate))
	eventStore := NewXyprotoEStore(creator, NewEncoder(eventFactory), NewDecoder(eventFactory))

	events1 := []es.Event{
		&mockEvent1{es.WithGuid{guid1}, "event1"},
		&mockEvent2{es.WithGuid{guid1}, "event2"},
		&mockEvent2{es.WithGuid{guid1}, "event4"},
		&mockEvent1{es.WithGuid{guid1}, "event6"},
	}
	events2 := []es.Event{
		&mockEvent1{es.WithGuid{guid2}, "event3"},
		&mockEvent1{es.WithGuid{guid2}, "event5"},
		&mockEvent1{es.WithGuid{guid2}, "event7"},
		&mockEvent2{es.WithGuid{guid2}, "event8"},
	}

	err := eventStore.Update(guid1, 0, events1)
	assert.NoError(t, err, "更新guid1事件不应该发生错误")
	err = eventStore.Update(guid2, 0, events2)
	assert.NoError(t, err, "更新guid2事件不应该发生错误")

	out1, version1 := eventStore.Find(guid1)
	out2, version2 := eventStore.Find(guid2)
	assert.Equal(t, len(events1), version1, "guid1返回的version有误")
	assert.Equal(t, len(events2), version2, "guid2返回的version有误")
	assert.Equal(t, events1, out1, "guid1返回的events有误")
	assert.Equal(t, events2, out2, "guid2返回的events有误")
}

func TestStoreNoRegisteredEvent(t *testing.T) {
	db, _ := simplebolt.New(path.Join(os.TempDir(), "bolt.db"))
	defer db.Close()
	guid := es.NewGuid()
	creator := simplebolt.NewCreator(db)
	defer func() {
		list, _ := creator.NewList("Event" + string(guid))
		list.Clear()
		list.Remove()
	}()
	eventFactory := es.NewEventFactory()
	eventFactory.RegisterAggregate(new(mockAggregate))
	eventStore := NewXyprotoEStore(creator, NewEncoder(eventFactory), NewDecoder(eventFactory))

	events := []es.Event{
		&mockEvent1{es.WithGuid{guid}, "event1"},
		&mockEvent2{es.WithGuid{guid}, "event2"},
		&mockNotRegisterEvent{es.WithGuid{guid}, "NotRegisterEvent"},
	}

	err := eventStore.Update(guid, 0, events)
	assert.NoError(t, err, "添加未注册的事件现在应该不报错")
	assert.Panics(t, func() { eventStore.Find(guid) }, "读取事件流中包含未注册的事件应该抛出异常")
}
