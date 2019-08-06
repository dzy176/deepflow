package mapreduce

import (
	"testing"

	"gitlab.x.lan/yunshan/droplet-libs/app"
	"gitlab.x.lan/yunshan/droplet-libs/codec"
	datatype "gitlab.x.lan/yunshan/droplet-libs/zerodoc"
)

func TestStash(t *testing.T) {
	f := datatype.Field{
		IP:       0x0a2102c8,
		IP1:      0x0a2102c9,
		L3EpcID:  2,
		L3EpcID1: 3,
		ACLGID:   10,

		TAPType: datatype.ToR,
	}
	meter := datatype.UsageMeter{
		SumPacketTx: 1,
		SumPacketRx: 2,
		SumBitTx:    4,
		SumBitRx:    5,
	}
	tag1 := f.NewTag(datatype.IPPath | datatype.L3EpcIDPath | datatype.Protocol | datatype.ServerPort | datatype.TAPType)
	meter1 := &datatype.UsageMeter{}
	*meter1 = meter
	doc1 := &app.Document{Timestamp: 0x12345678, Tag: tag1, Meter: meter1}

	tag2 := f.NewTag(datatype.ACLGID | datatype.ACLDirection | datatype.Direction | datatype.TAPType | datatype.IP)
	meter2 := &datatype.UsageMeter{}
	*meter2 = meter
	doc2 := &app.Document{Timestamp: 0x12345678, Tag: tag2, Meter: meter2}

	tag3 := f.NewTag(datatype.IPPath | datatype.L3EpcIDPath | datatype.Protocol | datatype.ServerPort | datatype.TAPType)
	meter3 := &datatype.UsageMeter{}
	*meter3 = meter
	doc3 := &app.Document{Timestamp: 0x12345678, Tag: tag3, Meter: meter3}

	tag4 := f.NewTag(datatype.ACLGID | datatype.ACLDirection | datatype.Direction | datatype.TAPType | datatype.IP)
	meter4 := &datatype.UsageMeter{}
	*meter4 = meter
	doc4 := &app.Document{Timestamp: 0x12345679, Tag: tag4, Meter: meter4}

	stash := NewSlidingStash(100, 30, 0)
	stash.Add([]interface{}{doc1, doc2, doc3, doc4})
	docs := stash.Dump()
	if len(docs) != 3 {
		t.Error("文档数量不正确")
	}

	e := &codec.SimpleEncoder{}
	if docs[0].(*app.Document).Tag.(*datatype.Tag).GetID(e) != tag1.GetID(e) {
		t.Error("文档0的tag不正确")
	}
	if docs[1].(*app.Document).Tag.(*datatype.Tag).GetID(e) != tag2.GetID(e) {
		t.Error("文档1的tag不正确")
	}
	if docs[2].(*app.Document).Tag.(*datatype.Tag).GetID(e) != tag4.GetID(e) {
		t.Error("文档2的tag不正确")
	}
	if docs[0].(*app.Document).Meter.(*datatype.UsageMeter).SumPacketTx != 2*meter.SumPacketTx {
		t.Error("文档0的meter不正确")
	}

	stash.Clear()
	stash.Add([]interface{}{doc1, doc2, doc3, doc4})
	if len(docs) != 3 {
		t.Error("Clear后文档数量不正确")
	}
}
