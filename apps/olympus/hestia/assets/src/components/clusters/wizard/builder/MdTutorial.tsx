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
    `1. All your apps are built from these building blocks. 

type TopologyBaseInfraWorkload struct {
   *v1core.Service       \`json:"service"\`
   *v1core.ConfigMap     \`json:"configMap"\`
   *v1.Deployment        \`json:"deployment"\`
   *v1.StatefulSet       \`json:"statefulSet"\`
   *v1networking.Ingress \`json:"ingress"\`
}

2. A Cluster is how you can link related building blocks.

    Containers -> Pods -> StatefulSet/Deployment -> Workload Base
           ConfigMap, Service, Ingress, .. etc  -> 

    1..N Workload Bases -> One Cluster Base
 
Workload bases that fall under one Cluster are only registered to this cluster base definition
 ->Each workload base can be switched out interchangeably. Allowing multiple configuration options
    Lighthouse:
        -> Goerli version
            -> v3.5.0
            -> v3.4.0
        -> Archive version
            -> 1 Ti storage
        -> etc..

Here's an example to illustrate this.
One cluster base can be called a consensus client:
        -> Lighthouse 
        -> Prysm
        -> etc

1..N Cluster Bases  -> One Cluster

You can define your beacon cluster.
    which contains:
        -> consensusClients
            -> Lighthouse, Prysm, Lodestar, etc
        -> execClients
            -> Geth, Erigon, etc..

You can then stack clusters with other clusters

Ethereum Beacon Cluster + Ethereum Validator Client Cluster = Staking Configuration 

Your apps are all compatible with any Cloud Provider, On-Premise, Hybrid, etc. You can route your app deploys
using this routing structure.

type CloudCtxNs struct {
    CloudProvider string \`json:"cloudProvider"\`
    Region        string \`json:"region"\`
    Context       string \`json:"context"\`
    Namespace     string \`json:"namespace"\`
    Env           string \`json:"env"\`
}
`