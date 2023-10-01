import React, {useEffect, useRef, useState} from "react";
import {SelectionText} from "@uiw/react-textarea-code-editor";
import MonacoEditor from "react-monaco-editor/lib/editor";
import {editor} from "monaco-editor";
import setTheme = editor.setTheme;

export default function MdTextField(props: any) {
    const [code, setCode] = useState('');

    const onChange = (textInput: string) => {
        setCode(textInput);
    };

    const themeRef = useRef<string>()
    function onSelectThemeChange(e: React.ChangeEvent<HTMLSelectElement>) {
        e.persist();
        document.documentElement.setAttribute('data-color-mode', /^vs$/.test(e.target.value) ? 'light' : 'dark');
        themeRef.current = e.target.value;
        setTheme(e.target.value);
    }
    useEffect(() => {
        if (themeRef.current) {
            // @ts-ignore
            const obj = new SelectionText(themeRef.current);
        }
    }, []);
    return (
        <div>
            <MonacoEditor
                height="800px"
                width="1000px"
                language="go"
                theme={'vs-dark'}
                onChange={(event) => onChange(event)}
                value={tutorial}
            />
        </div>
    );
}

const tutorial =
    `All your apps are built from these building blocks. 

type TopologyBaseInfraWorkload struct {
   *v1core.Service       \`json:"service"\`
   *v1core.ConfigMap     \`json:"configMap"\`
   *v1.Deployment        \`json:"deployment"\`
   *v1.StatefulSet       \`json:"statefulSet"\`
   *v1networking.Ingress \`json:"ingress"\`
}

A Cluster is how you can link related building blocks.

type CloudCtxNs struct {
    CloudProvider string \`json:"cloudProvider"\`
    Region        string \`json:"region"\`
    Context       string \`json:"context"\`
    Namespace     string \`json:"namespace"\`
    Env           string \`json:"env"\`
}

You route your app deploys using this routing structure.
`