import React from 'react';
import type {Recipe} from "./services/models/recipe.ts";
import {Button, Grid, Paper} from "@mui/material";
import {css} from "@emotion/react";
import LaunchIcon from '@mui/icons-material/Launch';
import CopyIcon from '@mui/icons-material/ContentCopyRounded'
const containerCss = css`
    overflow: hidden;
    width: 100%;
    height: 100%;
`;

const gridCss = css`
    margin: 1em auto 0.8em;
    padding: 0;
    height: fit-content;
    width: 70%;
`;

export const RecipeLinks:React.FC<{recipe: Recipe}> = ({recipe}) => {
    return <Paper elevation={3} css={containerCss}>
        <Grid container justifyContent="space-between" css={gridCss}>
            <Grid>
                <Button variant="contained" endIcon={<LaunchIcon/>} onClick={
                    ()=>window.open(`https://recipes.guardianapis.com/api/content/by-uid/${recipe.id}`, '_blank')
                }>Raw JSON</Button>
            </Grid>
            <Grid>
                <Button variant="contained" endIcon={<LaunchIcon/>} onClick={
                    ()=>window.open(`https://www.theguardian.com/${recipe.canonicalArticle}`)
                }>Website version</Button>
            </Grid>
            <Grid>
                <Button variant="contained" endIcon={<LaunchIcon/>} onClick={
                    ()=>window.open(`https://composer.gutools.co.uk/content/${recipe.composerId}`)
                }>Composer</Button>
            </Grid>
            <Grid>
                <Button variant="contained" endIcon={<CopyIcon/>} onClick={
                    ()=>navigator.clipboard.writeText(recipe.id)
                }>Copy ID</Button>
            </Grid>
        </Grid>
    </Paper>
}