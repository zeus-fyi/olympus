import {useMemo} from "react";
import {ScaleOrdinal} from "d3";

type AxisLeftProps = {
    yScale: ScaleOrdinal<string, any>;
    pixelsPerTick: number;
    xOffset?: number;
};

// tick length
const TICK_LENGTH = 6;  // Add your desired tick length value here

export const AxisLeft = ({ yScale, pixelsPerTick, xOffset = 0 }: AxisLeftProps) => {

    const ticks = useMemo(() => {
        return yScale.domain().map((value, index) => ({
            value,
            yOffset: yScale(value) + index + pixelsPerTick / 2,
        }));
    }, [yScale]);

    return (
        <g transform={`translate(${xOffset}, 0)`}>
            {/* Draw axis line */}
            <line x1="0" x2="0" y1="0" y2={yScale.range()[yScale.range().length - 1]} stroke="black" />

            {/* Draw ticks */}
            {ticks.map((tick, index) => (
                <g key={index} transform={`translate(0, ${tick.yOffset})`}>
                    <line x1="0" x2={-TICK_LENGTH} y1="0" y2="0" stroke="black" />
                    <text x={-TICK_LENGTH - 4} y={5} textAnchor="end">
                        {tick.value}
                    </text>
                </g>
            ))}
        </g>
    );
};

export default AxisLeft;
