import React, {useEffect, useState} from "react";
import {
    Alert,
    Grid,
    LinearProgress,
    ListItemButton, MenuItem,
    Pagination,
    Paper, Select,
    Snackbar,
    Stack,
    Typography
} from "@mui/material";
import {css} from "@emotion/react";
import List from '@mui/material/List';
import ListItemText from '@mui/material/ListItemText';
import {QueryReports} from "./services/DirectElasticLookup.ts";
import type {SingleHitResponse} from "./services/models/elastic.ts";
import type {Recipe} from "./services/models/recipe.ts";
import {cacheEverythingInList, retrieveCacheAsRecord} from "./services/RecipeCache.ts";
import type {CheckerOutput} from "./services/models/recipe-health-checker.ts";
import {z} from "zod";
import {FormatBestDate, FormatDate} from "./utils.ts";

const boundingCss = css`
    width: 95%;
    height: 95%;
    padding: 0.6em;
`;

const recipeList = css`
    height: 98%;
    overflow-y: scroll;
    overflow-x: hidden;
`;

interface RecipeListProps {
    onReportSelected: (rpt: SingleHitResponse)=>void;
    onRecipeLoaded: (recep: Recipe|undefined)=>void;
}

export const RecipeList:React.FC<RecipeListProps> = ({onReportSelected, onRecipeLoaded})=>{
    const [selectedIndex, setSelectedIndex] = React.useState(0);
    const [pageStart, setPageStart] = useState(0);
    const [pageSize, setPageSize] = useState(25);
    const [reports, setReports] = useState<SingleHitResponse[]>([]);
    const [loading, setLoading] = useState(false);
    const [lastError, setLastError] = useState<string|undefined>();
    const [recipeContent, setRecipeContent] = useState<Record<string, Recipe>>({});
    const [pageCount, setPageCount] = useState(1);

    useEffect(() => {
        setLoading(true);
        QueryReports(pageStart, pageSize)
            .then(
                (result)=>{
                    setReports(result.hits.hits);
                    setPageCount(result.hits.total.value ? Math.ceil(result.hits.total.value / pageSize) : 0);
                    cacheEverythingInList(result.hits.hits.map(r=>r._source.recipe_id))
                        .then(()=>{
                            setLoading(false);
                            setLastError(undefined);
                            setRecipeContent(retrieveCacheAsRecord());
                        })
                        .catch((err)=>{
                            setLoading(false);
                            if(err instanceof Error) {
                                setLastError(err.message)
                            } else {
                                const msg = String(err);
                                setLastError(msg);
                            }
                        })
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

    useEffect(() => {
        const rpt = reports[selectedIndex];

        if(rpt && recipeContent[rpt._source.recipe_id]) {
            onRecipeLoaded(recipeContent[rpt._source.recipe_id]);
        }
    }, [recipeContent, selectedIndex]);

    const handleListItemClick = (index: number) => {
        setSelectedIndex(index);
        onReportSelected(reports[index]);
    }

    const ImprovedListItem:React.FC<{report:CheckerOutput, source:z.infer<typeof Recipe>}> = ({report, source}) => <ListItemText
        primary={`${source.title}`}
        secondary={<Stack>
            <Typography variant="caption">{source.contributors?.join(",") ?? ""} {source.byline?.join(",") ?? ""}</Typography>
            <Typography variant="caption">Published {FormatBestDate(source)} / Tested {FormatDate(report.timestamp)}</Typography>
            <Typography variant="caption">{report.annotation_count} annotations</Typography>
        </Stack>}
        />;

    return <Paper elevation={3} css={boundingCss}>
        {
            loading ? <LinearProgress/> : undefined
        }
        <Grid container direction="row" justifyContent="space-between">
            <Grid>
            <Pagination count={pageCount}
                        size="large"
                        onChange={(_,pageNum)=>setPageStart((pageNum-1)*pageSize)}
            />
            </Grid>
            <Grid>
                <Select variant="filled"
                        label="Page size"
                        style={{width: '100px'}}
                        value={pageSize}
                        onChange={(evt)=>setPageSize(evt.target.value)}>
                    <MenuItem value={10}>10</MenuItem>
                    <MenuItem value={25}>25</MenuItem>
                    <MenuItem value={50}>50</MenuItem>
                    <MenuItem value={100}>100</MenuItem>
                </Select>
            </Grid>
        </Grid>
        <List css={recipeList}>
            {
                reports.map((rpt, idx)=>
                    <ListItemButton selected={selectedIndex===idx} onClick={()=>handleListItemClick(idx)} key={idx}>
                        {
                            recipeContent[rpt._source.recipe_id] ?
                                <ImprovedListItem report={rpt._source} source={recipeContent[rpt._source.recipe_id]}/> :
                                <ListItemText primary={rpt._source.recipe_id} secondary={`${rpt._source.annotation_count} annotations`}/>
                        }
                    </ListItemButton>
                )
            }
        </List>
        {
            lastError ?
                <Snackbar open={true}
                          anchorOrigin={{ vertical: 'top', horizontal: 'left'}}
                          autoHideDuration={5000}
                          onClose={()=>setLastError(undefined)}
                          >
                    <Alert severity="error">{lastError}</Alert>
                </Snackbar> : undefined
        }
    </Paper>
}