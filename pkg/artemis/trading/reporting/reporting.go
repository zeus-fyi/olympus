package artemis_reporting

/*
SELECT count(*), amount_out_addr, expected_profit_amount_out
FROM eth_mev_tx_analysis
WHERE end_reason = 'success' AND rx_block_number > 17639300 AND amount_in_addr = '0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2'
GROUP BY amount_out_addr, expected_profit_amount_out
*/
