package until

type ObjStruct interface {
	OutPrint()
}

func (s PHTR2) OutPrint() {
}

type Ns struct {
}

func (c Ns) Work(o ObjStruct) {
	o.OutPrint()
}
