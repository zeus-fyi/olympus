import {useMemo} from "react";
import {ScaleLinear} from "d3";

type AxisBottomProps = {
    xScale: ScaleLinear<number, number>;
};

// tick length
const TICK_LENGTH = 6;

export const AxisBottom = ({ xScale }: AxisBottomProps) => {
    const [min, max] = xScale.range();
    // Compute ticks: increase the number to have more ticks
    const numTicks = 13;
    const ticks = useMemo(() => {
        return xScale.ticks(numTicks).map(value => ({
            value,
            xOffset: xScale(value)
        }));
    }, [xScale, numTicks]);

    return (
        <>
            {/* Main horizontal line */}
            <path
                d={["M", min + 20, 0, "L", max - 20, 0].join(" ")}
                fill="none"
                stroke="currentColor"
            />
            {/* Ticks and labels */}
            {ticks.map(({ value, xOffset }) => (
                <g key={value} transform={`translate(${xOffset}, 0)`}>
                    <text
                        key={value}
                        style={{
                            fontSize: "20px",
                            textAnchor: "middle",
                            transform: "translateY(20px)",
                        }}
                    >
                        {value}
                    </text>
                </g>
            ))}
        </>
    );
};
