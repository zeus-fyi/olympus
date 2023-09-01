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
    return (
        <>
            {/* Vertical line */}
            <line
                x1={median}
                x2={median}
                y1={offset}
                y2={offset}
                stroke={stroke}
                strokeWidth={STROKE_WIDTH} // Corrected attribute name
            />
            {/* Rectangle box */}
            <rect
                x={q1}
                y={offset-50}
                width={q3-q1}
                height={100}
                stroke={stroke}
                fill={fill}
            />

            {/* horizontal line */}
            <line
                x1={min}
                x2={max}
                y1={offset}
                y2={offset}
                stroke={stroke}
                strokeWidth={STROKE_WIDTH} // Corrected attribute name
            />
            <g>
                <text
                    x={q1-18}
                    y={offset-18}  // Position below the line
                    style={{
                        fontSize: "12px",
                        textAnchor: "middle", // To center the text relative to given 'x'
                        dominantBaseline: "middle"
                    }}
                >
                    p25
                </text>
            </g>
            <g>
                <text
                    x={q3+18}
                    y={offset-18}  // Position below the line
                    style={{
                        fontSize: "12px",
                        textAnchor: "middle", // To center the text relative to given 'x'
                        dominantBaseline: "middle"
                    }}
                >
                    p75
                </text>
            </g>
            {/* Median line */}
            <g>
            <line
                x1={median}
                x2={median}
                y1={offset-100}
                y2={offset+200}
                stroke={stroke}
                strokeWidth={STROKE_WIDTH} // Corrected attribute name
            />
                {/* Label for the line */}
                <text
                    x={median}
                    y={offset-115}  // Position below the line
                    style={{
                        fontSize: "12px",
                        textAnchor: "middle", // To center the text relative to given 'x'
                        dominantBaseline: "middle"
                    }}
                >
                    median
                </text>
            </g>
        </>
    );
};
