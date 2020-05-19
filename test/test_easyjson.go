// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package test

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson785d9294DecodeGithubComUnravelinNullTest(in *jlexer.Lexer, out *pltest) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "A":
			(out.A).UnmarshalEasyJSON(in)
		case "B":
			(out.B).UnmarshalEasyJSON(in)
		case "C":
			(out.C).UnmarshalEasyJSON(in)
		case "D":
			(out.D).UnmarshalEasyJSON(in)
		case "E":
			(out.E).UnmarshalEasyJSON(in)
		case "A1":
			(out.A1).UnmarshalEasyJSON(in)
		case "B1":
			(out.B1).UnmarshalEasyJSON(in)
		case "C1":
			(out.C1).UnmarshalEasyJSON(in)
		case "D1":
			(out.D1).UnmarshalEasyJSON(in)
		case "E1":
			(out.E1).UnmarshalEasyJSON(in)
		case "A2":
			(out.A2).UnmarshalEasyJSON(in)
		case "B2":
			(out.B2).UnmarshalEasyJSON(in)
		case "C2":
			(out.C2).UnmarshalEasyJSON(in)
		case "D2":
			(out.D2).UnmarshalEasyJSON(in)
		case "E2":
			(out.E2).UnmarshalEasyJSON(in)
		case "A3":
			(out.A3).UnmarshalEasyJSON(in)
		case "B3":
			(out.B3).UnmarshalEasyJSON(in)
		case "C3":
			(out.C3).UnmarshalEasyJSON(in)
		case "D3":
			(out.D3).UnmarshalEasyJSON(in)
		case "E3":
			(out.E3).UnmarshalEasyJSON(in)
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson785d9294EncodeGithubComUnravelinNullTest(out *jwriter.Writer, in pltest) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"A\":"
		out.RawString(prefix[1:])
		(in.A).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"B\":"
		out.RawString(prefix)
		(in.B).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"C\":"
		out.RawString(prefix)
		(in.C).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"D\":"
		out.RawString(prefix)
		(in.D).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"E\":"
		out.RawString(prefix)
		(in.E).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"A1\":"
		out.RawString(prefix)
		(in.A1).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"B1\":"
		out.RawString(prefix)
		(in.B1).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"C1\":"
		out.RawString(prefix)
		(in.C1).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"D1\":"
		out.RawString(prefix)
		(in.D1).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"E1\":"
		out.RawString(prefix)
		(in.E1).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"A2\":"
		out.RawString(prefix)
		(in.A2).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"B2\":"
		out.RawString(prefix)
		(in.B2).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"C2\":"
		out.RawString(prefix)
		(in.C2).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"D2\":"
		out.RawString(prefix)
		(in.D2).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"E2\":"
		out.RawString(prefix)
		(in.E2).MarshalEasyJSON(out)
	}
	if (in.A3).IsDefined() {
		const prefix string = ",\"A3\":"
		out.RawString(prefix)
		(in.A3).MarshalEasyJSON(out)
	}
	if (in.B3).IsDefined() {
		const prefix string = ",\"B3\":"
		out.RawString(prefix)
		(in.B3).MarshalEasyJSON(out)
	}
	if (in.C3).IsDefined() {
		const prefix string = ",\"C3\":"
		out.RawString(prefix)
		(in.C3).MarshalEasyJSON(out)
	}
	if (in.D3).IsDefined() {
		const prefix string = ",\"D3\":"
		out.RawString(prefix)
		(in.D3).MarshalEasyJSON(out)
	}
	if (in.E3).IsDefined() {
		const prefix string = ",\"E3\":"
		out.RawString(prefix)
		(in.E3).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v pltest) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson785d9294EncodeGithubComUnravelinNullTest(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *pltest) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson785d9294DecodeGithubComUnravelinNullTest(l, v)
}
