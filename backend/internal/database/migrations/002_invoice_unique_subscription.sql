-- Ensure one invoice per subscription
ALTER TABLE invoices
	ADD UNIQUE KEY uq_invoice_subscription (subscription_id);
