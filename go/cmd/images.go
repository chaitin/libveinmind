package cmd

import (
	"context"

	"github.com/chaitin/libveinmind/go"
	"github.com/chaitin/libveinmind/go/plugin"
)

// ScanAllImages scans image provided by runtime list.
func ScanAllImages(
	ctx context.Context, rang plugin.ExecRange,
	runtime []api.Runtime, opts ...plugin.ExecOption,
) error {
	iter, err := plugin.IterateTyped(rang, "image")
	if err != nil {
		return err
	}
	return Scan(ctx, iter, runtime, opts...)
}

// ScanImages scans image provided by image list.
func ScanImages(
	ctx context.Context, rang plugin.ExecRange,
	images []api.Image, opts ...plugin.ExecOption,
) error {
	iter, err := plugin.IterateTyped(rang, "image")
	if err != nil {
		return err
	}
	return Scan(ctx, iter, images,
		plugin.WithPrependArgs("--id"),
		plugin.WithExecOptions(opts...))
}

// ScanImage scan an image provided.
func ScanImage(
	ctx context.Context, rang plugin.ExecRange,
	image api.Image, opts ...plugin.ExecOption,
) error {
	return ScanImages(ctx, rang, []api.Image{image}, opts...)
}

// ScanImageIDs with a runtime and a list of IDs provided.
func ScanImageIDs(
	ctx context.Context, rang plugin.ExecRange,
	runtime api.Runtime, ids []string, opts ...plugin.ExecOption,
) error {
	iter, err := plugin.IterateTyped(rang, "image")
	if err != nil {
		return err
	}
	return ScanIDs(ctx, iter, runtime, ids,
		plugin.WithPrependArgs("--id"),
		plugin.WithExecOptions(opts...))
}

// imageExactIDs specifies whether the argument list specifies
// ID instead of searchable names.
var imageExactIDs bool

// ImageIDsHandler is the handler for current list of images.
type ImageIDsHandler func(*Command, api.Runtime, []string) error

// MapImageIDsCommand attempts to create an image IDs command.
//
// The command will attempt to initialize the runtime object
// from specified mode with flags, scan and match images in
// the runtime, and collect those qualified image IDs.
func (idx *Index) MapImageIDsCommand(
	c *Command, f ImageIDsHandler,
) *Command {
	c = idx.MapModeCommand(c, "image", struct{}{}, func(
		c *Command, args []string, root interface{},
	) error {
		r, ok := root.(api.Runtime)
		if !ok {
			return IncompatibleMode()
		}
		var imageIDs []string
		if len(args) == 0 {
			ids, err := r.ListImageIDs()
			if err != nil {
				return err
			}
			imageIDs = append(imageIDs, ids...)
		} else if imageExactIDs {
			imageIDs = append(imageIDs, args...)
		} else {
			for _, arg := range args {
				ids, err := r.FindImageIDs(arg)
				if err != nil {
					return err
				}
				imageIDs = append(imageIDs, ids...)
			}
		}
		return f(c, r, imageIDs)
	})
	flags := c.PersistentFlags()
	flags.BoolVar(&imageExactIDs, "id", false,
		"whether fully qualified ID is specified")
	return c
}

// AddImageIDsCommand invokes MapImageCommand with no return.
func (idx *Index) AddImageIDsCommand(
	c *Command, f ImageIDsHandler,
) {
	_ = idx.MapImageIDsCommand(c, f)
}

// MapImageIDsCommand issues defaultIndex.MapImageIDsCommand.
func MapImageIDsCommand(
	c *Command, f ImageIDsHandler,
) *Command {
	return defaultIndex.MapImageIDsCommand(c, f)
}

// AddImageCommand issues defaultIndex.AddImageIDsCommand.
func AddImageIDsCommand(
	c *Command, f ImageIDsHandler,
) {
	defaultIndex.AddImageIDsCommand(c, f)
}

// ImageHandler is the handler for specified images.
type ImageHandler func(*Command, api.Image) error

// MapImageCommand attempts to create a image command.
//
// The command will attempt to initialize the runtime object
// from specified mode with flags, scan and match images in
// the runtime, and open matched images, one at once.
func (idx *Index) MapImageCommand(
	c *Command, f ImageHandler,
) *Command {
	return idx.MapImageIDsCommand(c, func(
		c *Command, r api.Runtime, imageIDs []string,
	) error {
		for _, imageID := range imageIDs {
			if err := func() error {
				image, err := r.OpenImageByID(imageID)
				if err != nil {
					return err
				}
				defer func() { _ = image.Close() }()
				return f(c, image)
			}(); err != nil {
				return err
			}
		}
		return nil
	})
}

// AddImageCommand invokes MapImageCommand with no return.
func (idx *Index) AddImageCommand(
	c *Command, f ImageHandler,
) {
	_ = idx.MapImageCommand(c, f)
}

// MapImageCommand issues defaultIndex.MapImageCommand.
func MapImageCommand(
	c *Command, f ImageHandler,
) *Command {
	return defaultIndex.MapImageCommand(c, f)
}

// AddImageCommand issues defaultIndex.AddImageCommand.
func AddImageCommand(
	c *Command, f ImageHandler,
) {
	defaultIndex.AddImageCommand(c, f)
}
