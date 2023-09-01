import {useMemo} from "react";
import {ScaleOrdinal} from "d3";

type AxisLeftProps = {
    yScale: ScaleOrdinal<string, any>;
    pixelsPerTick: number;
    height?: number;
};

export const AxisLeft = ({ yScale, pixelsPerTick, height = 0 }: AxisLeftProps) => {
    const ticks = useMemo(() => {
        return yScale.domain().map((value, index) => ({
            value,
            yOffset: yScale(value) ?? 0,
        }));
    }, [yScale]);

    // + index + pixelsPerTick / 2,
    return (
        <g transform={`translate(0, ${height})`}>
            {/* Draw axis line */}
            <line x1="0" x2="0" y1="0" y2={yScale.range()[yScale.range().length - 1]} stroke="black" />
            {/* Draw ticks */}
            {ticks.map((tick, index) => (
                <g key={index} transform={`translate(0, ${tick.yOffset})`}>
                    <line x1={-6} x2={4} y1="0" y2="0" stroke="black" />
                    <text x={-8} y={5} textAnchor="end">
                        {tick.value}
                    </text>
                </g>
            ))}
        </g>
    );
};

export default AxisLeft;
