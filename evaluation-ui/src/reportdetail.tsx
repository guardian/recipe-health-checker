import React from 'react';
import type {SingleHitResponse} from "./services/models/elastic.ts";
import {Paper} from "@mui/material";
import {css} from "@emotion/react";
import Markdown from "react-markdown";

interface ReportDetailProps {
    content: SingleHitResponse;
}

const boundingCss = css`
    padding: 0.2em;
    height: 100%;
    width: 100%;
`;

const scrollingMarkdown = css`
    overflow-y: scroll;
    height: 100%;
    width: 100%;
`;
export const ReportDetail:React.FC<ReportDetailProps> = ({content}) => {
    return <Paper css={boundingCss} elevation={3}>
        <div css={scrollingMarkdown}>
            <Markdown>
                {content._source.annotated_text}
            </Markdown>
        </div>
    </Paper>
}