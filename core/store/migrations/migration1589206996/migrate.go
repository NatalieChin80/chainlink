package migration1589206996

import (
	"github.com/jinzhu/gorm"
)

func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
	  	CREATE TABLE eth_transactions (
			id BIGSERIAL PRIMARY KEY,
			nonce bigint, 
			from_address bytea NOT NULL,
			to_address bytea NOT NULL,
			encoded_payload bytea NOT NULL,
			value numeric(78, 0) NOT NULL,
			created_at timestamptz NOT NULL
	  	);

		CREATE UNIQUE INDEX idx_eth_transactions_nonce_from_address ON eth_transactions (nonce, from_address);
		CREATE INDEX idx_eth_transactions_created_at ON eth_transactions USING BRIN (created_at);

		CREATE TABLE eth_transaction_attempts (
			id BIGSERIAL PRIMARY KEY,
		 	eth_transaction_id bigint REFERENCES eth_transactions (id) NOT NULL,
		 	gas_price numeric(78,0) NOT NULL,
		 	gas_limit bigint NOT NULL,
		 	signed_raw_tx bytea NOT NULL,
		 	hash bytea,
		 	error text,
		 	confirmed_in_block_num bigint,
		 	confirmed_in_block_hash bytea,
		 	confirmed_at timestamptz,
		 	created_at timestamptz NOT NULL
		);

		ALTER TABLE eth_transaction_attempts ADD CONSTRAINT chk_cannot_have_both_hash_and_error CHECK (
			(error IS NULL AND hash IS NULL AND confirmed_in_block_num IS NULL AND confirmed_in_block_hash IS NULL AND confirmed_at IS NULL)
			OR
			(error IS NOT NULL AND hash IS NULL AND confirmed_in_block_num IS NULL AND confirmed_in_block_hash IS NULL AND confirmed_at IS NULL)
			OR
			(error IS NULL AND hash IS NOT NULL AND confirmed_in_block_num IS NOT NULL AND confirmed_in_block_hash IS NOT NULL AND confirmed_at IS NOT NULL)
		);

		CREATE UNIQUE INDEX idx_eth_transaction_attempts_signed_raw_tx ON eth_transaction_attempts (signed_raw_tx);
		CREATE UNIQUE INDEX idx_eth_transaction_attempts_hash ON eth_transaction_attempts (hash);
		CREATE UNIQUE INDEX idx_eth_transaction_attempts_only_one_finalized ON eth_transaction_attempts (eth_transaction_id) WHERE (error IS NOT NULL) OR (hash IS NOT NULL);
		CREATE INDEX idx_eth_transaction_attempts ON eth_transaction_attempts (eth_transaction_id);
		CREATE INDEX idx_eth_transactions_confirmed_in_block_num ON eth_transaction_attempts (confirmed_in_block_num) WHERE confirmed_in_block_num IS NOT NULL;
		CREATE INDEX idx_eth_transactions_attempts_created_at ON eth_transaction_attempts USING BRIN (created_at);

		CREATE TABLE eth_task_run_transactions (
			task_run_id uuid NOT NULL REFERENCES task_runs (id) ON DELETE CASCADE,
			eth_transaction_id bigint NOT NULL REFERENCES eth_transactions (id) ON DELETE CASCADE
		);

		CREATE UNIQUE INDEX idx_eth_task_run_transactions_task_run_id ON eth_task_run_transactions (task_run_id);
		CREATE UNIQUE INDEX idx_eth_task_run_transactions_eth_transaction_id ON eth_task_run_transactions (eth_transaction_id);
	`).Error
}
