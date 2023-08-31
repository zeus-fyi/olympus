import {useMemo} from "react";
import {ScaleBand} from "d3";

type AxisLeftProps = {
    yScale: ScaleBand<string>;
    pixelsPerTick: number;
    xOffset?: number;
};

// tick length
const TICK_LENGTH = 6;

export const AxisLeft = ({ yScale, pixelsPerTick, xOffset = 0 as number }: AxisLeftProps) => {
    const [min, max] = yScale.range();

    const ticks = useMemo(() => {
        return yScale.domain().map((value) => ({
            value,
            // @ts-ignore
            yOffset: yScale(value) + yScale.bandwidth() / 2,
        }));
    }, [yScale]);

    return (
        <>
            {/* Main vertical line */}
            <line x1={xOffset} y1={0} x2={xOffset} y2={yScale.range()[1]} stroke="currentColor" />

            {/* Ticks and labels */}
            {ticks.map(({ value, yOffset }) => (
                <g key={value} transform={`translate(0, ${yOffset})`}>
                    <text
                        x={xOffset - TICK_LENGTH * 2}
                        y={0}
                        style={{
                            fontSize: "20px",
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
