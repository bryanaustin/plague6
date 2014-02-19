
package main

func DigestArgs(args []string) bool {
	var werr error

	if len(args) < 1 {
		Message("Supply at least one track file.")
		return true
	}
	if ako.Walks, werr = CompileWalkList(args...); werr != nil {
		Message("Error: %s", werr)
		return true
	}
	return false
}

func SetupDataAndTest() bool {
	ako.Data = make([]*WalkData, len(ako.Walks))
	for i := range ako.Walks {
		ako.Data[i] = new(WalkData)
		for _, steps := range ako.Walks[i].Steps {
			result := steps.Run()
			if result.Problem != nil {
				Message("Failed to run walk %d, error: %s", i, result)
				return true
			}
		}
	}
	return false
}