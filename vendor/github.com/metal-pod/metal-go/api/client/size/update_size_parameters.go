// Code generated by go-swagger; DO NOT EDIT.

package size

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/metal-pod/metal-go/api/models"
)

// NewUpdateSizeParams creates a new UpdateSizeParams object
// with the default values initialized.
func NewUpdateSizeParams() *UpdateSizeParams {
	var ()
	return &UpdateSizeParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewUpdateSizeParamsWithTimeout creates a new UpdateSizeParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewUpdateSizeParamsWithTimeout(timeout time.Duration) *UpdateSizeParams {
	var ()
	return &UpdateSizeParams{

		timeout: timeout,
	}
}

// NewUpdateSizeParamsWithContext creates a new UpdateSizeParams object
// with the default values initialized, and the ability to set a context for a request
func NewUpdateSizeParamsWithContext(ctx context.Context) *UpdateSizeParams {
	var ()
	return &UpdateSizeParams{

		Context: ctx,
	}
}

// NewUpdateSizeParamsWithHTTPClient creates a new UpdateSizeParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewUpdateSizeParamsWithHTTPClient(client *http.Client) *UpdateSizeParams {
	var ()
	return &UpdateSizeParams{
		HTTPClient: client,
	}
}

/*UpdateSizeParams contains all the parameters to send to the API endpoint
for the update size operation typically these are written to a http.Request
*/
type UpdateSizeParams struct {

	/*Body*/
	Body *models.V1SizeUpdateRequest

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the update size params
func (o *UpdateSizeParams) WithTimeout(timeout time.Duration) *UpdateSizeParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the update size params
func (o *UpdateSizeParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the update size params
func (o *UpdateSizeParams) WithContext(ctx context.Context) *UpdateSizeParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the update size params
func (o *UpdateSizeParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the update size params
func (o *UpdateSizeParams) WithHTTPClient(client *http.Client) *UpdateSizeParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the update size params
func (o *UpdateSizeParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the update size params
func (o *UpdateSizeParams) WithBody(body *models.V1SizeUpdateRequest) *UpdateSizeParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the update size params
func (o *UpdateSizeParams) SetBody(body *models.V1SizeUpdateRequest) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *UpdateSizeParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
