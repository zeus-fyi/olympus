import React, {useEffect, useRef, useState} from "react";
import {SelectionText} from "@uiw/react-textarea-code-editor";
import MonacoEditor from "react-monaco-editor/lib/editor";
import {editor} from "monaco-editor";
import {useSelector} from "react-redux";
import {RootState} from "../../../../redux/store";
import yaml from 'js-yaml';
import setTheme = editor.setTheme;

window.MonacoEnvironment = { getWorkerUrl: () => proxy };

let proxy = URL.createObjectURL(new Blob([`
	self.MonacoEnvironment = {
		baseUrl: 'https://unpkg.com/monaco-editor@latest/min/'
	};
	importScripts('https://unpkg.com/monaco-editor@latest/min/vs/base/worker/workerMain.js');
`], { type: 'text/javascript' }));

export const languageData = [
    'abap', 'aes', 'apex', 'azcli', 'bat', 'c', 'cameligo', 'clojure', 'coffeescript', 'cpp', 'csharp', 'csp', 'css', 'dart', 'dockerfile', 'fsharp', 'go', 'graphql', 'handlebars', 'hcl', 'html', 'ini', 'java', 'javascript', 'json', 'julia', 'kotlin', 'less', 'lex', 'lua', 'markdown', 'mips', 'msdax', 'mysql', 'objective', 'pascal', 'pascaligo', 'perl', 'pgsql', 'php', 'plaintext', 'postiats', 'powerquery', 'powershell', 'pug', 'python', 'r', 'razor', 'redis', 'redshift', 'restructuredtext', 'ruby', 'rust', 'sb', 'scala', 'scheme', 'scss', 'shell', 'sol', 'sql', 'st', 'swift', 'systemverilog', 'tcl', 'twig', 'typescript', 'vb', 'verilog', 'xml', 'yaml'
];

export default function YamlTextField(props: any) {
    const { previewType } = props;
    const clusterPreview = useSelector((state: RootState) => state.clusterBuilder.clusterPreview);
    const selectedComponentBaseName = useSelector((state: RootState) => state.clusterBuilder.selectedComponentBaseName);
    const selectedSkeletonBaseName = useSelector((state: RootState) => state.clusterBuilder.selectedSkeletonBaseName);
    const [code, setCode] = useState('');

    useEffect(() => {
        const clusterPreviewComponentBases = clusterPreview?.componentBases?.[selectedComponentBaseName];

        if (clusterPreviewComponentBases && Object.keys(clusterPreviewComponentBases).length > 0) {
            if (
                clusterPreviewComponentBases[selectedSkeletonBaseName] &&
                Object.keys(clusterPreviewComponentBases[selectedSkeletonBaseName]).length > 0
            ) {
                switch (previewType) {
                    case 'service':
                        setCode(yaml.dump(clusterPreviewComponentBases[selectedSkeletonBaseName].service));
                        break;
                    case 'configMap':
                        setCode(yaml.dump(clusterPreviewComponentBases[selectedSkeletonBaseName].configMap));
                        break;
                    case 'deployment':
                        setCode(yaml.dump(clusterPreviewComponentBases[selectedSkeletonBaseName].deployment));
                        break;
                    case 'statefulSet':
                        setCode(yaml.dump(clusterPreviewComponentBases[selectedSkeletonBaseName].statefulSet));
                        break;
                    case 'ingress':
                        setCode(yaml.dump(clusterPreviewComponentBases[selectedSkeletonBaseName].ingress));
                        break;
                    default:
                        break;
                }
            }
        }
    }, [previewType, clusterPreview, selectedComponentBaseName, selectedSkeletonBaseName]);

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
                    language="yaml"
                    theme={'vs-dark'}
                    onChange={(event) => onChange(event)}
                    value={code}
                />
            </div>
    );
}
