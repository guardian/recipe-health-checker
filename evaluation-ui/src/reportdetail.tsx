import React, {useEffect, useState} from 'react';
import type {SingleHitResponse} from "./services/models/elastic.ts";
import {
    Alert,
    CircularProgress,
    FormControlLabel,
    Paper,
    Snackbar,
    Switch,
    Typography
} from "@mui/material";
import {css} from "@emotion/react";
import Markdown from "react-markdown";
import {preProcessMarkdown, type PreProcessResults} from "./services/MDParser.ts";
import remarkBreaks from "remark-breaks";

interface ReportDetailProps {
    content: SingleHitResponse;
    selectedSection: string;
    showAnnotated: boolean;
}

const boundingCss = css`
    padding: 0.2em;
    height: 100%;
    width: 100%;
    overflow: hidden;
`;

const scrollingMarkdown = css`
    overflow-y: scroll;
    height: 100%;
    width: 100%;
`;

const progressSpinner = css`
    margin-left: auto;
    margin-right: auto;
`;
export const ReportDetail:React.FC<ReportDetailProps> = ({content, selectedSection, showAnnotated}) => {
    const [parsed, setParsed] = useState<PreProcessResults|undefined>();
    const [error, setError] = useState<string|undefined>();
    const [showAll, setShowAll] = useState(true);

    useEffect(() => {
        try {
            const result = preProcessMarkdown(content._source.annotated_text);
            console.log(result);

            setParsed(preProcessMarkdown(content._source.annotated_text));
        } catch(err) {
            console.error(err);
            if(err instanceof Error) {
                setError(err.message);
            } else {
                setError(String(err));
            }
        }
    }, [content]);

    useEffect(()=>{
        const entries = Array.from(parsed?.content.keys() ?? []);
        console.log(`${entries ? entries.join(";") : "no report detail content"}`, entries)
    }, [parsed]);

    useEffect(() => {
        console.log(`Section change: ${selectedSection}`);
    }, [selectedSection]);

    const RenderContent:React.FC = ()=>{
        if(showAll || selectedSection==="") {
            return <Markdown remarkPlugins={[remarkBreaks]}>
               {showAnnotated ? content._source.annotated_text : content._source.snapshot}
            </Markdown>
        }
        if(parsed?.content.has(selectedSection)) {
            console.log(parsed?.content.get(selectedSection));

            return <Markdown remarkPlugins={[remarkBreaks]}>{
                parsed?.content.get(selectedSection)?.map(b=>b.text).join("\n\n") ?? "(no content)"
            }</Markdown>
        } else {
            return <Typography variant="h5">Nothing present in that section</Typography>
        }
    }

    return <Paper css={boundingCss} elevation={3}>
                <FormControlLabel
                    control={<Switch checked={showAll} onChange={()=>setShowAll(prev=>!prev)}/>}
                    label="Show entire recipe"
                    style={{marginLeft: "auto"}}
                />
            <div css={scrollingMarkdown}>
                {
                    parsed ? <RenderContent/> : <CircularProgress css={progressSpinner}/>
                }
            </div>
        {
            error ? <Snackbar open={true} autoHideDuration={5000} onClose={()=>setError(undefined)}>
                <Alert severity="error">{error}</Alert>
            </Snackbar> : undefined
        }
    </Paper>
}