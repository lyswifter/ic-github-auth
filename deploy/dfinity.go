package deploy

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/lyswifter/ic-auth/types"
	"github.com/lyswifter/ic-auth/util"
)

func getController(targetpath string, islocal bool) (string, error) {
	// dfx wallet --network ic addresses

	var idcmd *exec.Cmd
	if islocal {
		idcmd = exec.Command("dfx", "wallet", "addresses")
	} else {
		idcmd = exec.Command("dfx", "wallet", "--network", "ic", "addresses")
	}

	idcmd.Dir = targetpath

	var b bytes.Buffer
	idcmd.Stderr = &b
	idcmd.Stdout = &b

	err := idcmd.Run()
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

func DeployWithDfx(targetpath string, f *os.File, repo string, islocal bool, framework string) ([]types.CanisterInfo, error) {

	var deploycmd *exec.Cmd
	if islocal {
		deploycmd = exec.Command("dfx", "deploy")
	} else {
		deploycmd = exec.Command("dfx", "deploy", "--network", "ic")
	}

	deploycmd.Dir = targetpath

	stderr, err := deploycmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	stdout, err := deploycmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	deploycmd.Start()

	errReader := bufio.NewReader(stderr)
	outReader := bufio.NewReader(stdout)

	canisterName := []string{}
	canisterId := []string{}

	for {
		line, err := errReader.ReadString('\n')
		if err == io.EOF {
			break
		}

		if err != nil {
			break
		}

		if strings.Contains(line, "canister_id") {
			name, id, err := extractCanisterInfo(line)
			if err != nil {
				break
			}

			canisterName = append(canisterName, name)
			canisterId = append(canisterId, id)
		}

		// write local
		_, err = f.WriteString(util.Format(line))
		if err != nil {
			break
		}
	}

	for {
		line, err := outReader.ReadString('\n')
		if err == io.EOF {
			break
		}

		if err != nil {
			break
		}

		if strings.Contains(line, "canister_id") {
			name, id, err := extractCanisterInfo(line)
			if err != nil {
				break
			}

			canisterName = append(canisterName, name)
			canisterId = append(canisterId, id)
		}

		// write local
		_, err = f.WriteString(util.Format(line))
		if err != nil {
			break
		}
	}

	type CanisterID struct {
		IC string `json:"ic"`
	}

	controller, err := getController(targetpath, islocal)
	if err != nil {
		return nil, err
	}

	var network string = "ic"
	if islocal {
		network = "local"
	}

	//read canister id
	cinfofile := filepath.Join(targetpath, "canister_ids.json")

	fmt.Printf("canister ids file path: %s\n", cinfofile)

	cinfos := []types.CanisterInfo{}

	if util.Exists(cinfofile) {
		ret, err := os.ReadFile(cinfofile)
		if err != nil {
			return nil, err
		}

		var infos map[string]CanisterID = make(map[string]CanisterID)
		err = json.Unmarshal(ret, &infos)
		if err != nil {
			return nil, err
		}

		fmt.Printf("canister info map: %+v", infos)

		for k, v := range infos {
			var ctype = "asssets"
			if framework == "dfx" && !strings.Contains(k, "assets") {
				ctype = "other"
			}

			cinfo := types.CanisterInfo{
				Repo:         repo,
				Controller:   controller,
				CanisterName: k,
				CanisterID:   v.IC,
				CanisterType: ctype,
				Framework:    framework,
				Network:      network,
			}
			cinfos = append(cinfos, cinfo)
		}
	} else {
		for i, v := range canisterName {
			id := canisterId[i]
			name := v

			cinfo := types.CanisterInfo{
				Repo:         repo,
				Controller:   controller,
				CanisterName: name,
				CanisterID:   id,
				CanisterType: "",
				Framework:    framework,
				Network:      network,
			}
			cinfos = append(cinfos, cinfo)
		}
	}

	deploycmd.Wait()

	return cinfos, nil
}

func extractCanisterInfo(input string) (string, string, error) {
	var split = "with canister_id"
	arr := strings.Split(input, split)

	first := strings.TrimSpace(arr[0])
	last := strings.TrimSpace(arr[1])

	first = strings.TrimSuffix(first, ",")
	firstarr := strings.Split(first, " ")
	canisterName := firstarr[len(firstarr)-1]

	fmt.Printf("input: %s canister name: %s canister id: %s", input, canisterName, last)
	return canisterName, last, nil
}
