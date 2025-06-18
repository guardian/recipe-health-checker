import React, {useState} from 'react';
import {css} from "@emotion/react";
import {Grid} from "@mui/material";
import {RecipeList} from "./recipelist.tsx";
import {SingleHitResponse} from "./services/models/elastic.ts";
import {ReportDetail} from "./reportdetail.tsx";
import {AnnotationDetails} from "./annotationdetails.tsx";

const boundingCss = css`
    width: 100vw;
    height: 100vh;
    margin: 0;
    padding: 1em;
    overflow: hidden;
    background-color: #c3c8e3;
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
    max-height: 47vh;
`;

export const MainWindow:React.FC = () => {
    //const routerParams = useParams();
    const [currentReport, setCurrentReport] = useState<SingleHitResponse|undefined>();
    const [selectedSection, setSelectedSection] = useState<string>(undefined);

    return <Grid container spacing={2} direction="row" css={boundingCss}>
            <Grid css={menuBox}>
                <RecipeList onReportSelected={setCurrentReport}/>
            </Grid>
            <Grid container spacing={2} direction="column" style={{overflow: 'hidden'}}>
                <Grid css={reportBox}>
                    {currentReport ? <AnnotationDetails report={currentReport._source} onSelectionChange={setSelectedSection}/> : undefined}
                </Grid>
                <Grid css={reportBox}>
                    {currentReport ? <ReportDetail content={currentReport} selectedSection={selectedSection} showAnnotated={false}/> : undefined}
                </Grid>
            </Grid>
        </Grid>
}