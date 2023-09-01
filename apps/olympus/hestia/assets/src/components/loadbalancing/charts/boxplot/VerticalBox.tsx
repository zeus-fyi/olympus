// A reusable component that builds a vertical box shape using svg
// Note: numbers here are px, not the real values in the dataset.

type VerticalBoxProps = {
    min: number;
    q1: number;
    median: number;
    q3: number;
    max: number;
    width: number;
    stroke: string;
    fill: string;
    offset: number;
};
const STROKE_WIDTH = 1; // Adjust this to your needs

export const VerticalBox = ({
                         min,
                         q1,
                         median,
                         q3,
                         max,
                         width,
                         stroke,
                         fill, offset,
                     }:VerticalBoxProps) => {

    let x = 0
    if (offset !== 0) {
        x = offset -1
    }
    console.log("offset", offset)
    return (
        <>
            {/* Vertical line */}
            <line
                x1={median}
                x2={median}
                y1={fill}
                y2={fill}
                stroke={stroke}
                strokeWidth={STROKE_WIDTH} // Corrected attribute name
            />
            {/* Rectangle box */}
            <rect
                x={x}
                y={offset}
                width={width}
                height={100}
                stroke={stroke}
                fill={fill}
            />
            {/* Median line */}
            <line
                x1={median}
                x2={median}
                y1={0}
                y2={offset}
                stroke={stroke}
                strokeWidth={STROKE_WIDTH} // Corrected attribute name
            />
        </>
    );
};
