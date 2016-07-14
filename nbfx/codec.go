package nbfx

import (
	"bytes"
)

type Encoder interface {
	Encode(*bytes.Reader) ([]byte, error)
}

type Decoder interface {
	Decode(*bytes.Reader) (string, error)
}

// for MultiByteInt31
const MASK_MBI31 = uint32(0x80) //0x80 = 128

const (
	EndElement byte = 0x01
	Comment    byte = 0x02
	Array      byte = 0x03

	ShortAttribute                byte = 0x04
	Attribute                     byte = 0x05
	ShortDictionaryAttribute      byte = 0x06
	DictionaryAttribute           byte = 0x07
	ShortXmlnsAttribute           byte = 0x08
	XmlnsAttribute                byte = 0x09
	ShortDictionaryXmlnsAttribute byte = 0x0A
	DictionaryXmlnsAttribute      byte = 0x0B
	PrefixDictionaryAttributeA    byte = 0x0C
	PrefixDictionaryAttributeB    byte = 0x0D
	PrefixDictionaryAttributeC    byte = 0x0E
	PrefixDictionaryAttributeD    byte = 0x0F
	PrefixDictionaryAttributeE    byte = 0x10
	PrefixDictionaryAttributeF    byte = 0x11
	PrefixDictionaryAttributeG    byte = 0x12
	PrefixDictionaryAttributeH    byte = 0x13
	PrefixDictionaryAttributeI    byte = 0x14
	PrefixDictionaryAttributeJ    byte = 0x15
	PrefixDictionaryAttributeK    byte = 0x16
	PrefixDictionaryAttributeL    byte = 0x17
	PrefixDictionaryAttributeM    byte = 0x18
	PrefixDictionaryAttributeN    byte = 0x19
	PrefixDictionaryAttributeO    byte = 0x1A
	PrefixDictionaryAttributeP    byte = 0x1B
	PrefixDictionaryAttributeQ    byte = 0x1C
	PrefixDictionaryAttributeR    byte = 0x1D
	PrefixDictionaryAttributeS    byte = 0x1E
	PrefixDictionaryAttributeT    byte = 0x1F
	PrefixDictionaryAttributeU    byte = 0x20
	PrefixDictionaryAttributeV    byte = 0x21
	PrefixDictionaryAttributeW    byte = 0x22
	PrefixDictionaryAttributeX    byte = 0x23
	PrefixDictionaryAttributeY    byte = 0x24
	PrefixDictionaryAttributeZ    byte = 0x25
	PrefixAttributeA              byte = 0x26
	PrefixAttributeB              byte = 0x27
	PrefixAttributeC              byte = 0x28
	PrefixAttributeD              byte = 0x29
	PrefixAttributeE              byte = 0x2A
	PrefixAttributeF              byte = 0x2B
	PrefixAttributeG              byte = 0x2C
	PrefixAttributeH              byte = 0x2D
	PrefixAttributeI              byte = 0x2E
	PrefixAttributeJ              byte = 0x2F
	PrefixAttributeK              byte = 0x30
	PrefixAttributeL              byte = 0x31
	PrefixAttributeM              byte = 0x32
	PrefixAttributeN              byte = 0x33
	PrefixAttributeO              byte = 0x34
	PrefixAttributeP              byte = 0x35
	PrefixAttributeQ              byte = 0x36
	PrefixAttributeR              byte = 0x37
	PrefixAttributeS              byte = 0x38
	PrefixAttributeT              byte = 0x39
	PrefixAttributeU              byte = 0x3A
	PrefixAttributeV              byte = 0x3B
	PrefixAttributeW              byte = 0x3C
	PrefixAttributeX              byte = 0x3D
	PrefixAttributeY              byte = 0x3E
	PrefixAttributeZ              byte = 0x3F

	ShortElement             byte = 0x40
	Element                  byte = 0x41
	ShortDictionaryElement   byte = 0x42
	DictionaryElement        byte = 0x43
	PrefixDictionaryElementA byte = 0x44
	PrefixDictionaryElementB byte = 0x45
	PrefixDictionaryElementC byte = 0x46
	PrefixDictionaryElementD byte = 0x47
	PrefixDictionaryElementE byte = 0x48
	PrefixDictionaryElementF byte = 0x49
	PrefixDictionaryElementG byte = 0x4A
	PrefixDictionaryElementH byte = 0x4B
	PrefixDictionaryElementI byte = 0x4C
	PrefixDictionaryElementJ byte = 0x4D
	PrefixDictionaryElementK byte = 0x4E
	PrefixDictionaryElementL byte = 0x4F
	PrefixDictionaryElementM byte = 0x50
	PrefixDictionaryElementN byte = 0x51
	PrefixDictionaryElementO byte = 0x52
	PrefixDictionaryElementP byte = 0x53
	PrefixDictionaryElementQ byte = 0x54
	PrefixDictionaryElementR byte = 0x55
	PrefixDictionaryElementS byte = 0x56
	PrefixDictionaryElementT byte = 0x57
	PrefixDictionaryElementU byte = 0x58
	PrefixDictionaryElementV byte = 0x59
	PrefixDictionaryElementW byte = 0x5A
	PrefixDictionaryElementX byte = 0x5B
	PrefixDictionaryElementY byte = 0x5C
	PrefixDictionaryElementZ byte = 0x5D
	PrefixElementA           byte = 0x5E
	PrefixElementB           byte = 0x5F
	PrefixElementC           byte = 0x60
	PrefixElementD           byte = 0x61
	PrefixElementE           byte = 0x62
	PrefixElementF           byte = 0x63
	PrefixElementG           byte = 0x64
	PrefixElementH           byte = 0x65
	PrefixElementI           byte = 0x66
	PrefixElementJ           byte = 0x67
	PrefixElementK           byte = 0x68
	PrefixElementL           byte = 0x69
	PrefixElementM           byte = 0x6A
	PrefixElementN           byte = 0x6B
	PrefixElementO           byte = 0x6C
	PrefixElementP           byte = 0x6D
	PrefixElementQ           byte = 0x6E
	PrefixElementR           byte = 0x6F
	PrefixElementS           byte = 0x70
	PrefixElementT           byte = 0x71
	PrefixElementU           byte = 0x72
	PrefixElementV           byte = 0x73
	PrefixElementW           byte = 0x74
	PrefixElementX           byte = 0x75
	PrefixElementY           byte = 0x76
	PrefixElementZ           byte = 0x77

	ZeroText                          byte = 0x80
	ZeroTextWithEndElement            byte = 0x81
	OneText                           byte = 0x82
	OneTextWithEndElement             byte = 0x83
	FalseText                         byte = 0x84
	FalseTextWithEndElement           byte = 0x85
	TrueText                          byte = 0x86
	TrueTextWithEndElement            byte = 0x87
	Int8Text                          byte = 0x88
	Int8TextWithEndElement            byte = 0x89
	Int16Text                         byte = 0x8A
	Int16TextWithEndElement           byte = 0x8B
	Int32Text                         byte = 0x8C
	Int32TextWithEndElement           byte = 0x8D
	Int64Text                         byte = 0x8E
	Int64TextWithEndElement           byte = 0x8F
	FloatText                         byte = 0x90
	FloatTextWithEndElement           byte = 0x91
	DoubleText                        byte = 0x92
	DoubleTextWithEndElement          byte = 0x93
	DecimalText                       byte = 0x94
	DecimalTextWithEndElement         byte = 0x95
	DateTimeText                      byte = 0x96
	DateTimeTextWithEndElement        byte = 0x97
	Chars8Text                        byte = 0x98
	Chars8TextWithEndElement          byte = 0x99
	Chars16Text                       byte = 0x9A
	Chars16TextWithEndElement         byte = 0x9B
	Chars32Text                       byte = 0x9C
	Chars32TextWithEndElement         byte = 0x9D
	Bytes8Text                        byte = 0x9E
	Bytes8TextWithEndElement          byte = 0x9F
	Bytes16Text                       byte = 0xA0
	Bytes16TextWithEndElement         byte = 0xA1
	Bytes32Text                       byte = 0xA2
	Bytes32TextWithEndElement         byte = 0xA3
	StartListText                     byte = 0xA4
	StartListTextWithEndElement       byte = 0xA5
	EndListText                       byte = 0xA6
	EndListTextWithEndElement         byte = 0xA7
	EmptyText                         byte = 0xA8
	EmptyTextWithEndElement           byte = 0xA9
	DictionaryText                    byte = 0xAA
	DictionaryTextWithEndElement      byte = 0xAB
	UniqueIdText                      byte = 0xAC
	UniqueIdTextWithEndElement        byte = 0xAD
	TimeSpanText                      byte = 0xAE
	TimeSpanTextWithEndElement        byte = 0xAF
	UuidText                          byte = 0xB0
	UuidTextWithEndElement            byte = 0xB1
	UInt64Text                        byte = 0xB2
	UInt64TextWithEndElement          byte = 0xB3
	BoolText                          byte = 0xB4
	BoolTextWithEndElement            byte = 0xB5
	UnicodeChars8Text                 byte = 0xB6
	UnicodeChars8TextWithEndElement   byte = 0xB7
	UnicodeChars16Text                byte = 0xB8
	UnicodeChars16TextWithEndElement  byte = 0xB9
	UnicodeChars32Text                byte = 0xBA
	UnicodeChars32TextWithEndElement  byte = 0xBB
	QNameDictionaryText               byte = 0xBC
	QNameDictionaryTextWithEndElement byte = 0xBD
)
