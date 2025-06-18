import type {CheckerOutput} from "./services/models/recipe-health-checker.ts";
import {Paper} from "@mui/material";
import {useEffect, useState} from "react";
import {type AnnotationSummary, GetAnnotationSummary} from "./services/AnnotationProcessor.ts";

interface AnnotationDetailsProps {
    report: CheckerOutput;
}

export const AnnotationDetails:React.FC<AnnotationDetailsProps> = ({report}) => {
    const [summary, setSummary] = useState<AnnotationSummary|undefined>();

    useEffect(() => {
        setSummary(GetAnnotationSummary(report));
    }, [report]);

    return <Paper elevation={3}>
        <h2>Checker findings</h2>
        <ul>
            {
                summary?.breakdown ?
                    Object.keys(summary.breakdown).map(k=><li>{k} - {summary.breakdown[k]} annotations</li>) :
                    undefined
            }
        </ul>
    </Paper>
}