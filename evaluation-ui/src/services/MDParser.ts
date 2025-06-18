import {remark} from 'remark';
import remarkBreaks from 'remark-breaks';
import type { Parent} from "mdast";

interface ContentBody {
    text: string;
    html: string;
    annotationCount: number;
}

export interface PreProcessResults {
    content: Map<string, ContentBody[]>;
}

function gatherTextChildren(node:Parent, joiner:string="\n") {
    return node.children
        .filter(n=>n.type==='text')
        .map(n=>n.value as string)
        .join(joiner);
}

const annotationMatcher = /<!HEALTH:(.*)>/;

function generateContentBody(src:Parent): ContentBody {
    const text = gatherTextChildren(src, "\n");
    const html = gatherTextChildren(src, "<br>");
    const matches = annotationMatcher.exec(text);
    const annotationCount = matches ? matches.length : 0;
    return {
        text,
        html,
        annotationCount
    }
}

export function preProcessMarkdown(md:string):PreProcessResults {
    const parsed = remark().use(remarkBreaks).parse(md);
    let blockAccumulator:Parent[] = [];
    let currentSection = "";
    const content:Map<string, ContentBody[]> = new Map();

    let firstSection = true;

    for (const c of parsed.children) {
        if(c.type==='heading') {    //treat headings as section markers
            if(blockAccumulator.length > 0) {
                const sectionName = firstSection ? "Title" : currentSection;
                content.set(sectionName, blockAccumulator.map(generateContentBody));
                blockAccumulator = [];
                firstSection = false;
            }
            currentSection = gatherTextChildren(c);

        } else {
            if((c as Parent).children) {
                blockAccumulator.push(c as Parent);
            }
        }
    }

    if(blockAccumulator.length > 0 && currentSection != "") {
        content.set(currentSection, blockAccumulator.map(generateContentBody));
    }

    return {
        content,
    }
}