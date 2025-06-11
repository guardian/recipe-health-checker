import React, {useEffect, useState} from "react";
import {ListItemButton, Paper, Typography} from "@mui/material";
import {css} from "@emotion/react";
import List from '@mui/material/List';
import ListItemText from '@mui/material/ListItemText';
import {QueryReports} from "./services/DirectElasticLookup.ts";
import type {SingleHitResponse} from "./services/models/elastic.ts";

const boundingCss = css`
    width: 95%;
    height: 95%;
    margin: 1em;
    padding: 0.6em;
`;

interface RecipeListProps {
    onReportSelected: (rpt: SingleHitResponse)=>void;
}

export const RecipeList:React.FC<RecipeListProps> = ({onReportSelected})=>{
    const [selectedIndex, setSelectedIndex] = React.useState(0);
    const [pageStart, setPageStart] = useState(0);
    const [pageSize, setPageSize] = useState(25);
    const [reports, setReports] = useState<SingleHitResponse[]>([]);
    const [loading, setLoading] = useState(false);
    const [lastError, setLastError] = useState<string|undefined>();

    useEffect(() => {
        setLoading(true);
        QueryReports(pageStart, pageSize)
            .then(
                (result)=>{
                    console.log(result);
                    setReports(result.hits.hits);
                    setLoading(false);
                    setLastError(undefined);
                }
            )
            .catch(
                (err)=>{
                    console.error(err);
                    setLoading(false);
                    if(err instanceof Error) {
                        setLastError(err.message);
                    } else {
                        const msg = String(err);
                        setLastError(msg);
                    }
                }
            )
    }, [pageStart, pageSize]);

    const handleListItemClick = (index: number) => {
        setSelectedIndex(index);
        onReportSelected(reports[index]);
    }

    return <Paper elevation={3} css={boundingCss}>
        <Typography>Found {reports.length} matching recipes</Typography>
        <List style={{overflow: "scroll", height: "100%"}}>
            {
                reports.map((rpt, idx)=>
                    <ListItemButton selected={selectedIndex===idx} onClick={()=>handleListItemClick(idx)} key={idx}>
                        <ListItemText primary={rpt._source.recipe_id} secondary={`${rpt._source.annotation_count} annotations`}/>
                    </ListItemButton>
                )
            }
        </List>
    </Paper>
}