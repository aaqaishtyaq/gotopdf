package topdf

import (
	"context"
	"fmt"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/devtool"
	"github.com/mafredri/cdp/protocol/network"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/protocol/target"
	"github.com/mafredri/cdp/rpcc"
	log "github.com/sirupsen/logrus"
)

type ChromeParameters struct {
	TargetURL       string
	Headers         map[string]string
	Orientation     string
	PrintBackground bool
	MarginTop       float64
	MarginRight     float64
	MarginBottom    float64
	MarginLeft      float64
}

// CreatePdf is the main method to create PDF using Chrome Dev Protocol given a URI
func CreatePdf(ctx context.Context, params ChromeParameters) ([]byte, error) {
	// Use the DevTools API to manage targets
	devt, err := devtool.New("http://127.0.0.1:9222").Version(ctx)
	if err != nil {
		return nil, err
	}

	// Open a new RPC connection to the Chrome Debugging Protocol target
	conn, err := rpcc.DialContext(ctx, devt.WebSocketDebuggerURL)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Create new browser context
	baseBrowser := cdp.NewClient(conn)
	newContextTarget, _ := baseBrowser.Target.CreateBrowserContext(ctx, target.NewCreateBrowserContextArgs())
	if err != nil {
		return nil, err
	}
	defer baseBrowser.Target.DisposeBrowserContext(ctx, target.NewDisposeBrowserContextArgs(newContextTarget.BrowserContextID))


	// Create a new blank target with the new browser context
	newTargetArgs := target.NewCreateTargetArgs("about:blank").SetBrowserContextID(newContextTarget.BrowserContextID)
	newTarget, err := baseBrowser.Target.CreateTarget(ctx, newTargetArgs)
	if err != nil {
		return nil, err
	}

	// Close the target when done
	// (In development, skip this step to leave tabs open!)
	closeTargetArgs := target.NewCloseTargetArgs(newTarget.TargetID)
	defer func() {
		closeReply, err := baseBrowser.Target.CloseTarget(ctx, closeTargetArgs)
		if err != nil || !closeReply.Success {
			log.Error(fmt.Sprintf("Could not close target for: %v because: %v", params.TargetURL, err))
		}
	}()

	// Connect the client to the new target
	newTargetWsURL := fmt.Sprintf("ws://127.0.0.1:9222/devtools/page/%s", newTarget.TargetID)
	newContextConn, err := rpcc.DialContext(ctx, newTargetWsURL)
	if err != nil {
		return nil, err
	}
	defer newContextConn.Close()
	c := cdp.NewClient(newContextConn)

	// Enable the runtime
	err = c.Runtime.Enable(ctx)
	if err != nil {
		return nil, err
	}

	// Enable the network
	err = c.Network.Enable(ctx, network.NewEnableArgs())
	if err != nil {
		return nil, err
	}
	// Enable events
	err = c.Page.Enable(ctx)
	if err != nil {
		return nil, err
	}
	// Create a client to listen for the load event to be fired
	loadEventFiredClient, _ := c.Page.LoadEventFired(ctx)
	defer loadEventFiredClient.Close()

	// Tell the page to navigate to the URL
	// urlParsed, _ := url.ParseRequestURI(urlRequest)
	navArgs := page.NewNavigateArgs(params.TargetURL)
	c.Page.Navigate(ctx, navArgs)

	// Wait for the page to finish loading
	loadEventFiredClient.Recv()

	// Print to PDF
	printToPDFArgs := page.NewPrintToPDFArgs().
		SetLandscape(params.Orientation == "Landscape").
		SetPrintBackground(params.PrintBackground).
		SetMarginTop(params.MarginTop).
		SetMarginRight(params.MarginRight).
		SetMarginBottom(params.MarginBottom).
		SetMarginLeft(params.MarginLeft)

	pdf, err := c.Page.PrintToPDF(ctx, printToPDFArgs)
	if err != nil{
		return nil, err
	}

	return pdf.Data, nil
}
