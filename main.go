package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"cuelang.org/go/cue"
	cueerrors "cuelang.org/go/cue/errors"
	cueformat "cuelang.org/go/cue/format"
	"cuelang.org/go/cue/load"
	cueflow "cuelang.org/go/tools/flow"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage: %s <path>\n", os.Args[0])
		os.Exit(1)
	}
	root, err := loadCue(os.Args[1])
	if err != nil {
		fmt.Println(cueerrors.Details(err, nil))
		os.Exit(1)
	}
	root, err = run(context.TODO(), root)
	if err != nil {
		fmt.Println(cueerrors.Details(err, nil))
		os.Exit(1)
	}
	err = root.Value().Validate(cue.Concrete(true))
	if err != nil {
		fmt.Println(cueerrors.Details(err, nil))
		os.Exit(1)
	}
	fmt.Printf("\n\n === END RESULT ===\n%s\n", cueDump(root.Value()))
	// fmt.Printf("%v\n", root.Value())
}

func loadCue(p string) (*cue.Instance, error) {
	cfg := &load.Config{}

	binsts := load.Instances([]string{p}, cfg)
	instances := cue.Build(binsts)
	root := cue.Merge(instances...)
	if root.Err != nil {
		return nil, root.Err
	}

	if err := root.Value().Validate(); err != nil {
		return nil, err
	}

	return root, nil
}

func run(ctx context.Context, root *cue.Instance) (*cue.Instance, error) {
	l := sync.Mutex{}

	taskFunc := func(v cue.Value) (cueflow.Runner, error) {
		if v.Kind() != cue.StructKind {
			return nil, nil
		}

		input := v.LookupPath(cue.ParsePath("input"))
		if !input.Exists() {
			return nil, nil
		}

		fmt.Printf("[detected task at %q]\n", v.Path().String())

		return cueflow.RunnerFunc(func(t *cueflow.Task) error {
			fmt.Printf("PROCESSING %q\n", t.Path().String())
			input := t.Value().LookupPath(cue.ParsePath("input"))
			if err := input.Err(); err != nil {
				return err
			}
			output, err := input.String()
			if err != nil {
				return err
			}

			return t.Fill(map[string]string{
				"output": fmt.Sprintf("from %s: %s", t.Path().String(), output),
			})
		}), nil
	}

	updateFunc := func(c *cueflow.Controller, t *cueflow.Task) error {
		if t == nil {
			return nil
		}

		if t.State() != cueflow.Terminated {
			return nil
		}

		var err error
		l.Lock()
		root, err = root.Fill(t.Value(), cuePathToStrings(t.Path())...)
		l.Unlock()
		if err != nil {
			return fmt.Errorf("filling %s: %w", t.Path().String(), err)
		}
		fmt.Printf("FILLED in %s: %s\n\n", t.Path().String(), cueDump(t.Value()))

		return nil
	}

	flow := cueflow.New(&cueflow.Config{
		UpdateFunc: updateFunc,
	}, root, taskFunc)

	fmt.Printf("TASKS\n")
	for _, t := range flow.Tasks() {
		deps := []string{}
		for _, d := range t.Dependencies() {
			deps = append(deps, d.Path().String())
		}
		fmt.Printf("  ===> %s: dependencies: %s\n", t.Path(), strings.Join(deps, ", "))
	}
	fmt.Printf("\n\n")

	return root, flow.Run(ctx)
}

func cuePathToStrings(p cue.Path) []string {
	selectors := p.Selectors()
	out := make([]string, len(selectors))
	for i, sel := range selectors {
		out[i] = sel.String()
	}
	return out
}

func cueDump(v cue.Value) string {
	src := v.Syntax(cue.Final())
	out, err := cueformat.Node(src)
	if err != nil {
		panic(err)
	}
	return string(out)
}
