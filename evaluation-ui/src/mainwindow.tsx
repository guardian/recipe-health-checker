import React, {useState} from 'react';
import {css} from "@emotion/react";
import {Grid, Paper} from "@mui/material";
import {RecipeList} from "./recipelist.tsx";
import {SingleHitResponse} from "./services/models/elastic.ts";
import {ReportDetail} from "./reportdetail.tsx";

const boundingCss = css`
    width: 100vw;
    height: 100vh;
    margin: 1em;
    overflow: hidden;
`;

const menuBox = css`
    width: 33vw;
    height: 100vh;
    overflow-x: hidden;
    overflow-y: auto;
`;

const reportBox = css`
    flex-grow: 1;
    margin: 0;
    max-width: 60vw;
    max-height: 50vh;
`;

export const MainWindow:React.FC = () => {
    //const routerParams = useParams();
    const [currentReport, setCurrentReport] = useState<SingleHitResponse|undefined>();


    return <div css={boundingCss}>
        <Grid container spacing={2} direction="row">
            <Grid css={menuBox}>
                <RecipeList onReportSelected={setCurrentReport}/>
            </Grid>
            <Grid container spacing={2} direction="column">
                <Grid css={reportBox}>
                    <Paper elevation={3} style={{width: "100%", height: "100%"}}></Paper>
                </Grid>
                <Grid css={reportBox}>
                    {currentReport ? <ReportDetail content={currentReport}/> : undefined}
                </Grid>
            </Grid>
        </Grid>
    </div>
}