import React from 'react';
import {useParams} from "react-router-dom";
import {css} from "@emotion/react";
import {Grid} from "@mui/material";
import {RecipeList} from "./recipelist.tsx";

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

export const MainWindow:React.FC = () => {
    //const routerParams = useParams();

    return <div css={boundingCss}>
        <Grid container>
            <Grid item css={menuBox}>
                <RecipeList/>
            </Grid>
        </Grid>
    </div>
}