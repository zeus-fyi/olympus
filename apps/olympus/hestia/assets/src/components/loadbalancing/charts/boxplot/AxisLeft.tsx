import {useMemo} from "react";
import {ScaleBand} from "d3";

type AxisLeftProps = {
    yScale: ScaleBand<string>;
    pixelsPerTick: number;
    xOffset?: number;
};

// tick length
const TICK_LENGTH = 6;  // Add your desired tick length value here

export const AxisLeft = ({ yScale, pixelsPerTick, xOffset = 0 }: AxisLeftProps) => {
    const [min, max] = yScale.range();

    const ticks = useMemo(() => {
        return yScale.domain().map((value, index) => ({
            value,
            index,  // Add index to be used as the current iteration number
            // @ts-ignore
            yOffset: yScale(value) + yScale.bandwidth() / 2,
        }));
    }, [yScale]);

    return (
        <>
            {/* Main vertical line */}
            <line x1={xOffset} y1={0} x2={xOffset} y2={max} stroke="currentColor" />

            {/* Ticks and labels */}
            {ticks.map(({ value, yOffset, index }) => (
                <g key={value} transform={`translate(0, ${yOffset})`}>
                    {/* Draw tick line */}
                    <line x1={xOffset - TICK_LENGTH} x2={xOffset} y1={index} y2={index} stroke="currentColor"/>

                    {/* Draw tick label */}
                    <text
                        x={xOffset - TICK_LENGTH * 2}
                        y={index}
                        style={{
                            fontSize: "16px",
                            textAnchor: "end",
                            dominantBaseline: "middle",
                        }}
                    >
                        {value}
                    </text>
                </g>
            ))}
        </>
    );
};

export default AxisLeft;
