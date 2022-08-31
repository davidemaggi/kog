package structs

type MergeResult struct {
	IsOk          bool
	DoneSomething bool
	Msg           string
	From          string
	To            string
	Details       []MergeResultEntry
}

type MergeResultEntry struct {
	Added    int
	Removed  int
	Modified int
	ObjType  string
}

func New_MergeResult() (ret MergeResult) {

	result := MergeResult{
		IsOk:          false,
		DoneSomething: false,
		Details:       make([]MergeResultEntry, 0),
	}

	result.Details = append(result.Details, New_MergeResultEntry(Context))
	result.Details = append(result.Details, New_MergeResultEntry(Cluster))
	result.Details = append(result.Details, New_MergeResultEntry(User))

	return result
}

func New_MergeResultEntry(ot string) (ret MergeResultEntry) {

	result := MergeResultEntry{
		Added:    0,
		Modified: 0,
		Removed:  0,
		ObjType:  ot,
	}

	return result
}

const (
	Context string = "context"
	Cluster        = "cluster"
	User           = "user"
)

const (
	Added    string = "added"
	Removed         = "removed"
	Modified        = "modified"
)

func AddAction(obj MergeResult, action string, objType string) MergeResult {

	for i, detail := range obj.Details {

		if detail.ObjType == objType {
			switch action {
			case Added:
				obj.Details[i].Added++
				break
			case Modified:
				obj.Details[i].Modified++
				break
			case Removed:
				obj.Details[i].Removed++
				break
			default:
				break
			}
		}

	}

	return obj
}
