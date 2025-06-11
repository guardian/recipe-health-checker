import React, {useState} from 'react';
import {css} from "@emotion/react";
import {Grid} from "@mui/material";
import {RecipeList} from "./recipelist.tsx";
import {SingleHitResponse} from "./services/models/elastic.ts";
import {ReportDetail} from "./reportdetail.tsx";

const boundingCss = css`
    width: 100vw;
    height: 100vh;
    margin: 0;
    overflow: hidden;
`;

const menuBox = css`
    width: 33vw;
    height: 100vh;
    margin: 0;
    overflow-x: hidden;
    overflow-y: auto;
`;

const reportBox = css`
    flex-grow: 1;
    max-width: 60vw;
`;

export const MainWindow:React.FC = () => {
    //const routerParams = useParams();
    const [currentReport, setCurrentReport] = useState<SingleHitResponse|undefined>();


    return <div css={boundingCss}>
        <Grid container spacing={2}>
            <Grid item css={menuBox}>
                <RecipeList onReportSelected={setCurrentReport}/>
            </Grid>
            <Grid item css={reportBox}>
                {currentReport ? <ReportDetail content={currentReport}/> : undefined}
            </Grid>
        </Grid>
    </div>
}