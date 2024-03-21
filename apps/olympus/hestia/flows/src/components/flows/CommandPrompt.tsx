import MonacoEditor from "react-monaco-editor/lib/editor";
import {useEffect, useRef} from "react";
import {editor} from "monaco-editor";
import setTheme = editor.setTheme;

export function MbTaskCmdPrompt(props: any) {
    const {code, setCode, height, width, language, onChange} = props;
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
                height={height}
                width={width}
                language={language}
                theme={'vs-dark'}
                onChange={(event) => onChange(event)}
                value={code}
                options={{
                    wordWrap: "on"
                }}
            />
        </div>
    );
}
