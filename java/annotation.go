package java

import "github.com/dhamidi/sai/classfile"

type Annotation struct {
	Type              string
	ElementValuePairs []ElementValuePair
}

type ElementValuePair struct {
	Name  string
	Value interface{}
}

func annotationsFromClassfile(anns []classfile.Annotation, cp classfile.ConstantPool) []Annotation {
	result := make([]Annotation, len(anns))
	for i, a := range anns {
		result[i] = Annotation{
			Type:              descriptorToTypeName(cp.GetUtf8(a.TypeIndex)),
			ElementValuePairs: elementValuePairsFromClassfile(a.ElementValuePairs, cp),
		}
	}
	return result
}

func elementValuePairsFromClassfile(pairs []classfile.ElementValuePair, cp classfile.ConstantPool) []ElementValuePair {
	result := make([]ElementValuePair, len(pairs))
	for i, p := range pairs {
		result[i] = ElementValuePair{
			Name:  cp.GetUtf8(p.ElementNameIndex),
			Value: elementValueToGo(p.Value, cp),
		}
	}
	return result
}

func elementValueToGo(ev classfile.ElementValue, cp classfile.ConstantPool) interface{} {
	switch ev.Tag {
	case 'B', 'C', 'I', 'S', 'Z':
		if idx, ok := ev.Value.(uint16); ok {
			if val, found := cp.GetInteger(idx); found {
				return val
			}
		}
	case 'D':
		if idx, ok := ev.Value.(uint16); ok {
			if val, found := cp.GetDouble(idx); found {
				return val
			}
		}
	case 'F':
		if idx, ok := ev.Value.(uint16); ok {
			if val, found := cp.GetFloat(idx); found {
				return val
			}
		}
	case 'J':
		if idx, ok := ev.Value.(uint16); ok {
			if val, found := cp.GetLong(idx); found {
				return val
			}
		}
	case 's':
		if idx, ok := ev.Value.(uint16); ok {
			return cp.GetUtf8(idx)
		}
	case 'e':
		if ecv, ok := ev.Value.(classfile.EnumConstValue); ok {
			return map[string]string{
				"type":  descriptorToTypeName(cp.GetUtf8(ecv.TypeNameIndex)),
				"value": cp.GetUtf8(ecv.ConstNameIndex),
			}
		}
	case 'c':
		if idx, ok := ev.Value.(uint16); ok {
			return descriptorToTypeName(cp.GetUtf8(idx))
		}
	case '@':
		if ann, ok := ev.Value.(classfile.Annotation); ok {
			return Annotation{
				Type:              descriptorToTypeName(cp.GetUtf8(ann.TypeIndex)),
				ElementValuePairs: elementValuePairsFromClassfile(ann.ElementValuePairs, cp),
			}
		}
	case '[':
		if arr, ok := ev.Value.(classfile.ArrayValue); ok {
			result := make([]interface{}, len(arr.Values))
			for i, v := range arr.Values {
				result[i] = elementValueToGo(v, cp)
			}
			return result
		}
	}
	return nil
}

func descriptorToTypeName(desc string) string {
	if len(desc) == 0 {
		return ""
	}
	if desc[0] == 'L' && desc[len(desc)-1] == ';' {
		return classfile.InternalToSourceName(desc[1 : len(desc)-1])
	}
	return desc
}
