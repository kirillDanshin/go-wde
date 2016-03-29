package wde

type Cursor int

const (
	NoneCursor Cursor = iota
	NormalCursor
	ResizeNCursor
	ResizeECursor
	ResizeSCursor
	ResizeWCursor
	ResizeEWCursor
	ResizeNSCursor
	ResizeNECursor
	ResizeSECursor
	ResizeSWCursor
	ResizeNWCursor
	CrosshairCursor
	IBeamCursor
	GrabHoverCursor
	GrabActiveCursor
	NotAllowedCursor
	customCursorBase // custom cursors are numbered starting here
)

type CursorCtl interface {
	Set(id Cursor)
	Hide()
	Show()
}

/* TODO: custom cursors: func CreateCursor(draw.Image, hotspot image.Point) Cursor

func (c Cursor) IsCustom() bool {
	return c >= customCursorBase
}
*/
