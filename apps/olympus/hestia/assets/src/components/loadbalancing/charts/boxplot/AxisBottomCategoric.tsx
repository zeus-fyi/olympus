import {useMemo} from "react";
import {ScaleBand} from "d3";

type AxisBottomProps = {
    xScale: ScaleBand<string>;
};

// tick length
const TICK_LENGTH = 6;

export const AxisBottom = ({ xScale }: AxisBottomProps) => {
    const [min, max] = xScale.range();

    const ticks = useMemo(() => {
        return xScale.domain().map((value) => ({
            value,
            // @ts-ignore
            xOffset: xScale(value) + xScale.bandwidth() / 2,
        }));
    }, [xScale]);

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
