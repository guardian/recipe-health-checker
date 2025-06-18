import type {CheckerOutput} from "./services/models/recipe-health-checker.ts";
import {Grid, List, ListItemButton, Paper, Typography} from "@mui/material";
import React, {useEffect, useState} from "react";
import {type Annotation, type AnnotationSummary, GetAnnotationSummary} from "./services/AnnotationProcessor.ts";
import {css} from "@emotion/react";
import ListItemText from "@mui/material/ListItemText";

interface AnnotationDetailsProps {
    report: CheckerOutput;
    onSelectionChange: (newval:string)=>void;
}

const containerStyle = css`
    width: 100%;
    height: 100%;
    padding: 0.8em;
    overflow-y: auto;
    overflow-x: hidden;
`;

const selectorList = css`
    max-width: 30%;
    min-width: 150px;
`;

export const AnnotationDetails:React.FC<AnnotationDetailsProps> = ({report, onSelectionChange}) => {
    const [summary, setSummary] = useState<AnnotationSummary|undefined>();
    const [selectedSection, setSelectedSection] = useState("");
    const [relevantAnnotations, setRelevantAnnotations] = useState<Annotation[]|undefined>();

    useEffect(() => {
        setSummary(GetAnnotationSummary(report));
    }, [report]);

    useEffect(()=>{
        setRelevantAnnotations(
            summary?.content.filter(ann=>ann.section===selectedSection)
        );
        onSelectionChange(selectedSection);
    }, [selectedSection, summary]);

    return <Paper elevation={3} css={containerStyle}>
        <Typography variant="h4">Checker findings</Typography>
        <Grid container direction="row" spacing={1}>
            <Grid css={selectorList}>
                <List>
                    {
                        summary?.breakdown ?
                            Object.keys(summary.breakdown).map((k, idx)=>
                                <ListItemButton key={idx}
                                                selected={k===selectedSection}
                                                onClick={()=>setSelectedSection(k)}
                                >
                                    <ListItemText primary={k} secondary={`${summary.breakdown[k]} annotations`}/>
                                </ListItemButton>) :
                            undefined
                    }
                </List>
            </Grid>
            <Grid style={{flex: 1}}>
                {
                    relevantAnnotations ? <List>
                        {
                            relevantAnnotations.map((ann, idx)=>
                                <ListItemText key={idx} primary={ann.detail} secondary={ann.section}/>
                            )
                        }
                    </List> : undefined
                }
            </Grid>
        </Grid>
    </Paper>
}