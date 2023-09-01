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
};
const STROKE_WIDTH = 2; // Adjust this to your needs

export const VerticalBox = ({
                         min,
                         q1,
                         median,
                         q3,
                         max,
                         width,
                         stroke,
                         fill,
                     }:VerticalBoxProps) => {
    return (
        <>
            {/* Vertical line */}
            <line
                x1={width / 2}
                x2={width / 2}
                y1={min}
                y2={max}
                stroke={stroke}
                strokeWidth={STROKE_WIDTH} // Corrected attribute name
            />
            {/* Rectangle box */}
            <rect
                x={0}
                y={min}
                width={width}
                height={q3 - q1}
                stroke={stroke}
                fill={fill}
            />
            {/* Median line */}
            <line
                x1={0}
                x2={width}
                y1={median}
                y2={median}
                stroke={stroke}
                strokeWidth={STROKE_WIDTH} // Corrected attribute name
            />
        </>
    );
};
