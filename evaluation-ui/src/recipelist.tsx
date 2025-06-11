import React from "react";
import {ListItemButton, Paper} from "@mui/material";
import {css} from "@emotion/react";
import List from '@mui/material/List';
import ListItemText from '@mui/material/ListItemText';

const boundingCss = css`
    width: 95%;
    height: 95%;
    margin: 1em;
`
export const RecipeList:React.FC = ()=>{
    const [selectedIndex, setSelectedIndex] = React.useState(0);

    const handleListItemClick = (index: number) => {
        setSelectedIndex(index);
    }

    return <Paper elevation={3} css={boundingCss}>
        <List sx={{bgcolor: 'background.paper'}}>
            <ListItemButton selected={selectedIndex===0} onClick={()=>handleListItemClick(0)}>
                <ListItemText primary="recipe one" secondary="summary here"/>
            </ListItemButton>
            <ListItemButton selected={selectedIndex===0} onClick={()=>handleListItemClick(0)}>
                <ListItemText primary="recipe two" secondary="summary here"/>
            </ListItemButton>
        </List>
    </Paper>
}